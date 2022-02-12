package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math"
	"net/http"
	"time"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	goNanoid "github.com/matoous/go-nanoid/v2"
)

/*
TODO:
- Earn via share links
- Admin API
- Live notifications of transaction

ADDED:
- Earn via referral/invitation links
- Earn via cheques (contests, etc.)
- Earn via last idle session computation
*/

type VirtualWallet struct {
	UserID  string  `db:"user_id" json:"-"`
	Balance float32 `db:"balance" json:"balance"`
}

type VirtualTransaction struct {
	ID          string    `db:"id" json:"id"`
	UserID      string    `db:"user_id" json:"user_id"`
	Amount      float32   `db:"amount" json:"amount"`
	Description string    `db:"description" json:"description"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}

type Cheque struct {
	ID          string    `db:"id" json:"id"`
	UserID      string    `db:"user_id" json:"user_id"`
	Amount      float32   `db:"amount" json:"amount"`
	Description string    `db:"description" json:"description"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}

const virtualTransactionsTableName = "virtual_transactions"
const virtualWalletsTableName = "virtual_wallets"
const chequesTableName = "cheques"

type VirtualBank struct {
	DB *sqlx.DB
}

var accountCreated = time.Date(2022, time.February, 6, 0, 0, 0, 0, time.UTC).UnixMilli()

func (b *VirtualBank) AddInitialAmountToExistingAccounts(firebaseApp *firebase.App) error {
	authClient, err := firebaseApp.Auth(context.Background())
	if err != nil {
		return err
	}

	rows, err := b.DB.Query("SELECT user_id from associated_ids")
	if err != nil {
		return err
	}
	defer rows.Close()

	identifiers := []auth.UserIdentifier{}
	for rows.Next() {
		var userId string
		rows.Scan(&userId)
		identifiers = append(identifiers, auth.UIDIdentifier{UID: userId})
	}

	iters := int(math.Ceil(float64(len(identifiers)) / 100))
	added := 0

	lastFirstNum := 0
	lastLastNum := 100

	Transact(b.DB, func(tx *sqlx.Tx) error {
		for i := 1; i <= iters; i++ {
			if idLen := len(identifiers); idLen < lastLastNum {
				lastLastNum = idLen
			}

			results, err := authClient.GetUsers(context.Background(), identifiers[lastFirstNum:lastLastNum-1])
			if err != nil {
				return err
			}

			for _, user := range results.Users {
				if user.UserMetadata.CreationTimestamp > accountCreated {
					continue
				}

				if b.AddInitialBalanceTo(user.UID, tx); err != nil {
					return err
				} else {
					added++
				}
			}

			lastFirstNum, lastLastNum = lastLastNum+1, 100*(i+1)
		}

		return nil
	})

	log.Printf("%d existing accounts were loaded with virtual coins\n", added)
	return nil
}

func (b *VirtualBank) AddInitialBalanceTo(uid string, tx *sqlx.Tx) (*VirtualTransaction, error) {
	amount := float32(4000)
	if vWallet, err := b.GetWalletByUID(uid); err == nil && vWallet.Balance <= 0 {
		return b.AddBalanceTo(uid, amount, "Add balance", tx)
	}

	gotTransaction, err := b.AddTransaction(uid, amount, "Initial Balance", tx)
	if err != nil {
		return nil, err
	}

	walletSql, walletArgs, _ := psql.Insert(virtualWalletsTableName).Columns("user_id", "balance").Values(uid, amount).ToSql()
	if res, err := tx.Exec(walletSql, walletArgs...); err != nil {
		return nil, err
	} else if err := wrapSqlResult(res); err != nil {
		return nil, err
	}

	return gotTransaction, nil
}

func (b *VirtualBank) GetWalletByUID(uid string) (*VirtualWallet, error) {
	vWallet := &VirtualWallet{}
	if err := b.DB.Get(
		vWallet,
		fmt.Sprintf("SELECT * FROM %s WHERE user_id = $1", virtualWalletsTableName),
		uid,
	); err != nil {
		return nil, err
	}
	return vWallet, nil
}

func (b *VirtualBank) AddBalanceTo(uid string, amount float32, desc string, tx *sqlx.Tx) (*VirtualTransaction, error) {
	vWallet, err := b.GetWalletByUID(uid)
	if err == sql.ErrNoRows {
		return b.AddInitialBalanceTo(uid, tx)
	} else if err != nil {
		return nil, err
	}

	gotTransaction, err := b.AddTransaction(uid, amount, desc, tx)
	if err != nil {
		return nil, err
	}

	walletSql, walletArgs, _ := psql.Update(virtualWalletsTableName).Set("balance", vWallet.Balance+amount).Where(sq.Eq{"user_id": uid}).ToSql()
	if res, err := tx.Exec(walletSql, walletArgs...); err != nil {
		return nil, err
	} else if err := wrapSqlResult(res); err != nil {
		return nil, err
	}

	return gotTransaction, nil
}

func (b *VirtualBank) AddTransaction(uid string, amount float32, desc string, tx *sqlx.Tx) (*VirtualTransaction, error) {
	id, _ := goNanoid.New()
	newTransaction := &VirtualTransaction{
		ID:          id,
		UserID:      uid,
		Amount:      amount,
		Description: desc,
		CreatedAt:   time.Now(),
	}

	transactionSql := fmt.Sprintf(
		"INSERT INTO %s (id, user_id, amount, description, created_at) VALUES (:id, :user_id, :amount, :description, :created_at)",
		virtualTransactionsTableName,
	)

	if res, err := tx.NamedExec(transactionSql, newTransaction); err != nil {
		return nil, err
	} else if err := wrapSqlResult(res); err != nil {
		return nil, err
	}

	return newTransaction, nil
}

func (b *VirtualBank) DeductBalanceTo(uid string, amount float32, desc string, tx *sqlx.Tx) (float32, error) {
	vWallet, err := b.GetWalletByUID(uid)
	if err != nil {
		return 0, &ResponseError{
			WError:     err,
			StatusCode: http.StatusForbidden,
		}
	} else if vWallet.Balance-amount <= 0 {
		return 0, &ResponseError{
			StatusCode: http.StatusBadRequest,
			Message:    "Insufficient balance.",
		}
	}

	if _, err := b.AddTransaction(uid, -amount, desc, tx); err != nil {
		return 0, err
	}

	walletSql, walletArgs, _ := psql.Update(virtualWalletsTableName).Set("balance", vWallet.Balance-amount).Where(sq.Eq{"user_id": uid}).ToSql()
	if res, err := tx.Exec(walletSql, walletArgs...); err != nil {
		return 0, err
	} else if err := wrapSqlResult(res); err != nil {
		return 0, err
	}

	return vWallet.Balance - amount, nil
}

func (b *VirtualBank) ConvertIdleTime(uid string, lastActiveAt time.Time) (*VirtualTransaction, error) {
	// compute by the minute
	// 20 * (total idle time / 10 minutes)
	idleTime := time.Since(lastActiveAt)
	quotientMinutes := math.Floor(idleTime.Minutes() / 10)

	// player should be atleast 10 minutes in order to earn
	if quotientMinutes < 10 {
		return nil, nil
	}

	var trans *VirtualTransaction
	coinsToEarn := 20 * quotientMinutes
	err := Transact(b.DB, func(tx *sqlx.Tx) error {
		gotTransaction, err := b.AddBalanceTo(uid, float32(coinsToEarn), "Idle time earning", tx)
		if gotTransaction != nil {
			trans = gotTransaction
		}
		return err
	})
	return trans, err
}

func (b *VirtualBank) GenerateCheque(toUID string, amount float32) (*Cheque, error) {
	id, _ := goNanoid.New()
	timestamp := time.Now()
	chequeDescription := "Deposit check"
	res, err := b.DB.Exec(
		"INSERT INTO cheques (id, user_id, amount, description, created_at) VALUES ($1, $2, $3, $4, $6)",
		id, toUID, amount, chequeDescription, timestamp,
	)
	if err != nil {
		return nil, err
	} else if err := wrapSqlResult(res); err != nil {
		return nil, err
	}

	finalCheque := &Cheque{
		ID:          id,
		UserID:      toUID,
		Amount:      amount,
		Description: chequeDescription,
		CreatedAt:   timestamp,
	}

	return finalCheque, nil
}

func (b *VirtualBank) GenerateChequeBySID(sid string, amount float32) (*Cheque, error) {
	gotUser, err := getAssociatedUserBy(b.DB, sq.Eq{"associated_id": sid})
	if err != nil {
		return nil, err
	}

	return b.GenerateCheque(gotUser.UserID, amount)
}

func (b *VirtualBank) DepositChequeByID(chequeID string, recipientUID string) error {
	if len(chequeID) == 0 {
		return &ResponseError{
			StatusCode: http.StatusBadRequest,
			Message:    "Cheque ID is empty.",
		}
	} else if len(recipientUID) == 0 {
		return &ResponseError{
			WError:     fmt.Errorf("recipient uid is empty"),
			StatusCode: http.StatusBadRequest,
		}
	}

	// verify cheque
	gotCheque := Cheque{}
	if err := b.DB.Get(
		&gotCheque,
		"SELECT * FROM cheques WHERE id = $1 AND user_id = $2",
		chequeID, recipientUID,
	); err != nil {
		if err == sql.ErrNoRows {
			return &ResponseError{
				StatusCode: http.StatusNotFound,
				Message:    "invalid cheque",
			}
		}
		return err
	}

	return Transact(b.DB, func(tx *sqlx.Tx) error {
		// deposit cheque
		_, err := b.AddBalanceTo(
			gotCheque.UserID,
			gotCheque.Amount,
			gotCheque.Description,
			tx,
		)
		return err
	})
}

type virtualWalletKey struct{}

func getVirtualWalletFromReq(r *http.Request) *VirtualWallet {
	vWallet, ok := r.Context().Value(virtualWalletKey{}).(*VirtualWallet)
	if !ok {
		return nil
	}
	return vWallet
}

func checkBalance(b *VirtualBank) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return wrapHandler(func(rw http.ResponseWriter, r *http.Request) error {
			token := getAuthTokenByReq(r)
			if token == nil {
				return &ResponseError{
					StatusCode: http.StatusForbidden,
					Message:    "Forbidden transaction.",
				}
			}

			vWallet, err := b.GetWalletByUID(token.UID)
			if err != nil {
				return err
			}

			if vWallet.Balance <= 0 {
				return &ResponseError{
					StatusCode: http.StatusBadRequest,
					Message:    "Account balance is zero.",
				}
			}

			ctx := context.WithValue(r.Context(), virtualWalletKey{}, vWallet)
			next.ServeHTTP(rw, r.WithContext(ctx))
			return nil
		})
	}
}
