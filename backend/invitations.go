package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	goNanoid "github.com/matoous/go-nanoid/v2"
)

/*

TODO:
- Invite by e-mail? and notify them?

*/

const invitationMoney float32 = 150.0

type UserInvitationCode struct {
	ID        string    `db:"id" json:"id"`
	UserID    string    `db:"user_id" json:"-"`
	MaxUsers  int       `db:"max_users" json:"max_users"`
	UserCount int       `db:"user_count" json:"user_count"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	ExpiresAt time.Time `db:"expires_at" json:"expires_at"`
}

type InvitationSystem struct {
	DB         *sqlx.DB
	CookieName string
}

func (sys *InvitationSystem) CheckEligibilityByUID(uid string) error {
	invLinkCounts := 0
	if rows, err := sys.DB.Query("SELECT count(*) FROM user_invitation_codes WHERE user_id = ?", uid); err == nil && rows.Next() {
		rows.Scan(&invLinkCounts)
		rows.Close()
	}

	if invLinkCounts >= 5 {
		return &ResponseError{
			StatusCode: http.StatusBadRequest,
			Message:    "You are only allowed to generate 5 invitation codes / links at any given moment.",
		}
	}

	return nil
}

const defaultInvitationExpiration = 2 * time.Hour

func (sys *InvitationSystem) Generate(uid string, maxUsers int) (string, error) {
	if maxUsers <= 0 && maxUsers > 5 {
		return "", &ResponseError{
			StatusCode: http.StatusBadRequest,
			Message:    "Max users should be minimum of 1 and maximum of 5.",
		}
	} else if len(uid) == 0 {
		return "", &ResponseError{
			StatusCode: http.StatusBadRequest,
			Message:    "User ID is zero.",
		}
	}

	newInvitationId, _ := goNanoid.New()
	newInvCode := UserInvitationCode{
		ID:        newInvitationId,
		UserID:    uid,
		MaxUsers:  maxUsers,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(defaultInvitationExpiration),
	}

	if res, err := sys.DB.NamedExec(
		"INSERT INTO user_invitation_codes (id, user_id, max_users, created_at, expires_at) VALUES (:id, :user_id, :max_users, :created_at, :expires_at)",
		&newInvCode,
	); err != nil {
		return "", err
	} else if err := wrapSqlResult(res); err != nil {
		return "", err
	}

	return newInvitationId, nil
}

func (sys *InvitationSystem) VerifyInvitationCode(invCode string) (*UserInvitationCode, error) {
	gotInvitation := &UserInvitationCode{}

	if err := sys.DB.Get(gotInvitation, "SELECT * FROM user_invitation_codes WHERE id = ?", invCode); err != nil {
		if err == sql.ErrNoRows {
			return nil, &ResponseError{
				StatusCode: http.StatusNotFound,
				Message:    "Invitation link provided is invalid or has expired.",
			}
		}

		return nil, err
	}

	if time.Now().After(gotInvitation.ExpiresAt) {
		// destroy invitation after expiration
		tx, err := sys.DB.BeginTxx(context.Background(), &sql.TxOptions{})
		if err == nil {
			if err := sys.DestroyInvitation(gotInvitation.ID, tx); err != nil {
				log.Println(err)
			}
		}

		return nil, &ResponseError{
			StatusCode: http.StatusNotFound,
			Message:    "Invitation link provided is invalid or has expired.",
		}
	}

	return gotInvitation, nil
}

func (sys *InvitationSystem) DestroyInvitation(id string, tx *sqlx.Tx) error {
	if res, err := tx.Exec("DELETE FROM user_invitation_Codes WHERE id = ?", id); err != nil {
		passivePrintError(tx.Rollback())
		return err
	} else if err := wrapSqlResult(res); err != nil {
		passivePrintError(tx.Rollback())
		return err
	}

	return nil
}

func passivePrintError(err error) {
	if err != nil {
		log.Println(err)
	}
}

func (sys *InvitationSystem) UseInvitationFromReq(rw http.ResponseWriter, r *http.Request, signedUpUid string, b *VirtualBank) error {
	inv, invErr := sys.InvitationFromCookie(r)
	if invErr != nil {
		return invErr
	}

	tx, err := sys.DB.BeginTxx(r.Context(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	// dispose invitation if new user count reaches max users
	if inv.UserCount+1 == inv.MaxUsers {
		// delete invitation
		if err := sys.DestroyInvitation(inv.ID, tx); err != nil {
			return err
		}
	} else {
		// update invitation
		if res, err := tx.Exec("UPDATE user_invitation_codes SET user_count = ? WHERE id = ?", inv.UserCount+1, inv.ID); err != nil {
			passivePrintError(tx.Rollback())
			return err
		} else if err := wrapSqlResult(res); err != nil {
			passivePrintError(tx.Rollback())
			return err
		}
	}

	// give money to origin user
	if err := b.AddBalanceTo(inv.UserID, invitationMoney, "Invitation incentive", tx); err != nil {
		passivePrintError(tx.Rollback())
		return err
	}

	// give money to newly-created user
	if err := b.AddBalanceTo(
		signedUpUid, invitationMoney,
		fmt.Sprintf("Accepted invitation from %s", signedUpUid), tx,
	); err != nil {
		passivePrintError(tx.Rollback())
		return err
	}

	if err := tx.Commit(); err != nil {
		log.Println(tx.Rollback())
		return err
	}

	// delete cookie after usage
	http.SetCookie(rw, &http.Cookie{
		Name:   sys.CookieName,
		Value:  inv.ID,
		MaxAge: -1,
		Path:   "/",
	})

	return nil
}

func (sys *InvitationSystem) InvitationFromCookie(r *http.Request) (*UserInvitationCode, error) {
	invCookie, err := r.Cookie(sys.CookieName)
	if err != nil {
		return nil, err
	}

	return sys.VerifyInvitationCode(invCookie.Value)
}

func (sys *InvitationSystem) InjectToRequest(rw http.ResponseWriter, inv *UserInvitationCode) {
	cookie := &http.Cookie{
		Name:    sys.CookieName,
		Value:   inv.ID,
		Expires: inv.ExpiresAt,
		Path:    "/",
	}

	http.SetCookie(rw, cookie)
}

func (sys *InvitationSystem) GetInvitationsByUID(uid string) ([]UserInvitationCode, error) {
	invitations := []UserInvitationCode{}
	if err := sys.DB.Select(&invitations, "SELECT * FROM user_invitation_codes WHERE user_id = ?", uid); err != nil {
		if err == sql.ErrNoRows {
			return invitations, nil
		}

		return nil, err
	}

	return invitations, nil
}
