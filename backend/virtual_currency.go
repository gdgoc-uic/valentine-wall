package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
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
- Earn via referral/invitation links
- Earn via share links
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

const virtualTransactionsTableName = "virtual_transactions"
const virtualWalletsTableName = "virtual_wallets"

type VirtualBank struct {
	DB *sqlx.DB
}

var accountCreated = time.Date(2022, time.February, 6, 0, 0, 0, 0, time.UTC)

func (b *VirtualBank) AddInitialAmountToExistingAccounts(firebaseApp *firebase.App) error {
	authClient, err := firebaseApp.Auth(context.Background())
	if err != nil {
		return err
	}

	res, err := authClient.GetUsers(context.Background(), []auth.UserIdentifier{})
	if err != nil {
		return err
	}

	tx, err := b.DB.BeginTxx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	added := 0
	for _, user := range res.Users {
		if time.UnixMilli(user.UserMetadata.CreationTimestamp).Before(accountCreated) {
			if b.AddInitialBalanceTo(user.UID, tx); err != nil {
				if err := tx.Rollback(); err != nil {
					log.Println(err)
				}
				return err
			} else {
				added++
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	log.Printf("%d existing accounts were loaded with virtual coins\n", added)
	return nil
}

func (b *VirtualBank) AddInitialBalanceTo(uid string, tx *sqlx.Tx) error {
	amount := float32(4000)
	if err := b.AddTransaction(uid, amount, "Initial Balance", tx); err != nil {
		return err
	}

	walletSql, walletArgs, _ := sq.Insert(virtualWalletsTableName).Columns("user_id", "balance").Values(uid, amount).ToSql()
	if res, err := tx.Exec(walletSql, walletArgs...); err != nil {
		return err
	} else if err := wrapSqlResult(res); err != nil {
		return err
	}

	return nil
}

func (b *VirtualBank) GetWalletByUID(uid string) (*VirtualWallet, error) {
	vWallet := &VirtualWallet{}
	if err := b.DB.Get(
		vWallet,
		fmt.Sprintf("SELECT * FROM %s WHERE user_id = ?", virtualWalletsTableName),
		uid,
	); err != nil {
		return nil, err
	}
	return vWallet, nil
}

func (b *VirtualBank) AddBalanceTo(uid string, amount float32, tx *sqlx.Tx) error {
	vWallet, err := b.GetWalletByUID(uid)
	if err == sql.ErrNoRows {
		return b.AddInitialBalanceTo(uid, tx)
	} else if err != nil {
		return err
	}

	if err := b.AddTransaction(uid, amount, "Add balance", tx); err != nil {
		return err
	}

	walletSql, walletArgs, _ := sq.Update(virtualWalletsTableName).Set("balance", vWallet.Balance+amount).Where(sq.Eq{"user_id": uid}).ToSql()
	if res, err := tx.Exec(walletSql, walletArgs...); err != nil {
		return err
	} else if err := wrapSqlResult(res); err != nil {
		return err
	}

	return nil
}

func (b *VirtualBank) AddTransaction(uid string, amount float32, desc string, tx *sqlx.Tx) error {
	id, _ := goNanoid.New()
	transactionSql, txArgs, _ := sq.Insert(virtualTransactionsTableName).
		Columns("id", "user_id", "amount", "description").Values(id, uid, amount, desc).ToSql()
	if res, err := tx.Exec(transactionSql, txArgs...); err != nil {
		return err
	} else if err := wrapSqlResult(res); err != nil {
		return err
	}

	return nil
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

	if err := b.AddTransaction(uid, amount, desc, tx); err != nil {
		return 0, err
	}

	walletSql, walletArgs, _ := sq.Update(virtualWalletsTableName).Set("balance", vWallet.Balance-amount).Where(sq.Eq{"user_id": uid}).ToSql()
	if res, err := tx.Exec(walletSql, walletArgs...); err != nil {
		return 0, err
	} else if err := wrapSqlResult(res); err != nil {
		return 0, err
	}

	return vWallet.Balance - amount, nil
}

func checkBalance(b *VirtualBank) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return wrapHandler(func(rw http.ResponseWriter, r *http.Request) error {
			token, ok := r.Context().Value("authToken").(*auth.Token)
			if !ok {
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

			ctx := context.WithValue(r.Context(), "virtualWallet", vWallet)
			next.ServeHTTP(rw, r.WithContext(ctx))
			return nil
		})
	}
}
