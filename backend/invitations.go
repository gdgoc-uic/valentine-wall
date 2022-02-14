package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
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
	if rows, err := sys.DB.Query("SELECT count(*) FROM user_invitation_codes WHERE user_id = $1", uid); err == nil && rows.Next() {
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

	newInvCode := UserInvitationCode{
		UserID:    uid,
		MaxUsers:  maxUsers,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(defaultInvitationExpiration),
	}

	if rows, err := sys.DB.NamedQuery(
		"INSERT INTO user_invitation_codes (user_id, max_users, created_at, expires_at) VALUES (:user_id, :max_users, :created_at, :expires_at) RETURNING id",
		&newInvCode,
	); err != nil {
		return "", err
	} else if newInvitationId, err := wrapSqlRowsAfterInsert(rows); err != nil {
		return "", err
	} else {
		return newInvitationId, nil
	}
}

func (sys *InvitationSystem) VerifyInvitationCode(invCode string) (*UserInvitationCode, error) {
	gotInvitation := &UserInvitationCode{}

	if err := sys.DB.Get(gotInvitation, "SELECT * FROM user_invitation_codes WHERE id = $1", invCode); err != nil {
		if err == sql.ErrNoRows {
			return nil, &ResponseError{
				StatusCode: http.StatusNotFound,
				Message:    "Invitation link provided is invalid or has expired.",
			}
		}

		return nil, err
	}

	if time.Now().After(gotInvitation.ExpiresAt) {
		err := Transact(sys.DB, func(tx *sqlx.Tx) error {
			// destroy invitation after expiration
			if err := sys.DestroyInvitation(gotInvitation.ID, tx); err != nil {
				return err
			}
			return nil
		})

		if err == nil {
			return nil, &ResponseError{
				StatusCode: http.StatusNotFound,
				Message:    "Invitation link provided is invalid or has expired.",
			}
		}
	}

	return gotInvitation, nil
}

func (sys *InvitationSystem) DestroyInvitation(id string, tx *sqlx.Tx) error {
	if res, err := tx.Exec("DELETE FROM user_invitation_Codes WHERE id = $1", id); err != nil {
		return err
	} else if err := wrapSqlResult(res); err != nil {
		return err
	}
	return nil
}

func (sys *InvitationSystem) UseInvitationFromReq(rw http.ResponseWriter, r *http.Request, signedUpUid string, b *VirtualBank) error {
	inv, invErr := sys.InvitationFromCookie(r)
	if invErr != nil {
		return invErr
	}

	// dispose invitation if new user count reaches max users
	err := Transact(sys.DB, func(tx *sqlx.Tx) error {
		if inv.UserCount+1 == inv.MaxUsers {
			// delete invitation
			return sys.DestroyInvitation(inv.ID, tx)
		} else {
			// update invitation
			if res, err := tx.Exec("UPDATE user_invitation_codes SET user_count = $1 WHERE id = $2", inv.UserCount+1, inv.ID); err != nil {
				return err
			} else if err := wrapSqlResult(res); err != nil {
				return err
			}
		}

		// give money to origin user
		if invWallet, err := b.GetWalletByUID(inv.UserID); err == nil {
			if _, _, err := b.AddBalanceTo(invWallet, invitationMoney, "Invitation incentive", tx); err != nil {
				return err
			}
		} else {
			return err
		}

		// give money to newly-created user
		if newWallet, err := b.GetWalletByUID(signedUpUid); err == nil {
			if _, _, err := b.AddBalanceTo(
				newWallet, invitationMoney,
				fmt.Sprintf("Accepted invitation from %s", signedUpUid), tx,
			); err != nil {
				return err
			}
		} else {
			return err
		}

		return nil
	})

	if err == nil {
		// delete cookie after usage
		http.SetCookie(rw, &http.Cookie{
			Name:   sys.CookieName,
			Value:  inv.ID,
			MaxAge: -1,
			Path:   "/",
		})
	}

	return err
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
	if err := sys.DB.Select(&invitations, "SELECT * FROM user_invitation_codes WHERE user_id = $1", uid); err != nil {
		if err == sql.ErrNoRows {
			return invitations, nil
		}

		return nil, err
	}

	return invitations, nil
}
