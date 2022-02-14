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
	ID      string  `db:"id" json:"-"`
	UserID  string  `db:"user_id" json:"-"`
	Balance float32 `db:"balance" json:"balance"`
}

func (vw *VirtualWallet) Update(vTx *VirtualTransaction, tx *sqlx.Tx) (float32, error) {
	rows, err := tx.Query(
		fmt.Sprintf(
			"UPDATE %s SET balance = $1 WHERE user_id = $2 RETURNING balance",
			virtualWalletsTableName,
		),
		vw.Balance+vTx.Amount,
		vw.UserID,
	)
	defer rows.Close()
	if err != nil {
		return 0, err
	}

	var newBalance float32
	if !rows.Next() {
		errMessage := "Unable to process your transaction. Please try again."
		return 0, &ResponseError{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    errMessage,
		}
	}

	rows.Scan(&newBalance)
	vw.Balance = newBalance
	return newBalance, nil
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

	rows, err := b.DB.Query("SELECT user_id from associated_users")
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

func (b *VirtualBank) AddInitialBalanceTo(uid string, tx *sqlx.Tx) (*VirtualTransaction, float32, error) {
	amount := float32(4000)
	walletSql := fmt.Sprintf("INSERT INTO %s (user_id, balance) VALUES ($1, $2) RETURNING *", virtualWalletsTableName)

	rows, err := tx.Queryx(walletSql, uid, amount)
	if err != nil {
		return nil, 0, err
	}

	if !rows.Next() {
		errMessage := "Unable to create wallet. Please try again."
		return nil, 0, &ResponseError{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    errMessage,
		}
	}

	newWallet := &VirtualWallet{}
	rows.StructScan(newWallet)
	rows.Close()
	return b.AddBalanceTo(newWallet, amount, "Initial Balance", tx)
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

func (b *VirtualBank) AddBalanceTo(wallet *VirtualWallet, amount float32, desc string, tx *sqlx.Tx) (*VirtualTransaction, float32, error) {
	gotTransaction, err := b.AddTransaction(wallet, amount, desc, tx)
	if err != nil {
		return nil, wallet.Balance, err
	} else if newBalance, err := wallet.Update(gotTransaction, tx); err != nil {
		return nil, wallet.Balance, err
	} else {
		return gotTransaction, newBalance, nil
	}
}

func (b *VirtualBank) AddTransaction(wallet *VirtualWallet, amount float32, desc string, tx *sqlx.Tx) (*VirtualTransaction, error) {
	if wallet == nil {
		return nil, &ResponseError{
			StatusCode: http.StatusForbidden,
		}
	} else if wallet.Balance+amount <= 0 {
		return nil, &ResponseError{
			StatusCode: http.StatusBadRequest,
			WError:     fmt.Errorf("(account=%s, balance=%.2f, amount=%.2f, amountAfter=%.2f)", wallet.UserID, wallet.Balance, amount, wallet.Balance+amount),
			Message:    "Insufficient balance.",
		}
	}

	newTransaction := &VirtualTransaction{
		UserID:      wallet.UserID,
		Amount:      amount,
		Description: desc,
		CreatedAt:   time.Now(),
	}

	transactionSql := fmt.Sprintf(
		"INSERT INTO %s (user_id, amount, description, created_at) VALUES (:user_id, :amount, :description, :created_at) RETURNING id",
		virtualTransactionsTableName,
	)

	if rows, err := tx.NamedQuery(transactionSql, newTransaction); err != nil {
		return nil, err
	} else if id, err := wrapSqlRowsAfterInsert(rows); err != nil {
		return nil, err
	} else {
		newTransaction.ID = id
	}

	return newTransaction, nil
}

func (b *VirtualBank) DeductBalanceTo(vWallet *VirtualWallet, amount float32, desc string, tx *sqlx.Tx) (*VirtualTransaction, float32, error) {
	if vWallet == nil {
		return nil, 0, fmt.Errorf("empty wallet")
	}

	gotTransaction, err := b.AddTransaction(vWallet, -amount, desc, tx)
	if err != nil {
		return nil, vWallet.Balance, err
	} else if newBalance, err := vWallet.Update(gotTransaction, tx); err != nil {
		return nil, vWallet.Balance, err
	} else {
		return gotTransaction, newBalance, nil
	}
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
	wallet, err := b.GetWalletByUID(uid)
	if err != nil {
		return nil, err
	}

	if err := Transact(b.DB, func(tx *sqlx.Tx) error {
		gotTransaction, _, err := b.AddBalanceTo(wallet, float32(coinsToEarn), "Idle time earning", tx)
		if gotTransaction != nil {
			trans = gotTransaction
		}
		return err
	}); err != nil {
		return nil, err
	}

	return trans, nil
}

func (b *VirtualBank) GenerateCheque(toUID string, amount float32) (*Cheque, error) {
	timestamp := time.Now()
	chequeDescription := "Deposit check"
	var id string
	if rows, err := b.DB.Queryx(
		"INSERT INTO cheques (user_id, amount, description, created_at) VALUES ($1, $2, $3, $4, $6) RETURNING id",
		toUID, amount, chequeDescription, timestamp,
	); err != nil {
		return nil, err
	} else if id, err = wrapSqlRowsAfterInsert(rows); err != nil {
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

	recipientWallet, err := b.GetWalletByUID(recipientUID)
	if err != nil {
		return err
	}

	return Transact(b.DB, func(tx *sqlx.Tx) error {
		// deposit cheque
		_, _, err := b.AddBalanceTo(
			recipientWallet,
			gotCheque.Amount,
			gotCheque.Description,
			tx,
		)
		return err
	})
}

type virtualWalletKey struct{}

func injectWallet(b *VirtualBank) func(http.Handler) http.Handler {
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

			ctx := context.WithValue(r.Context(), virtualWalletKey{}, vWallet)
			next.ServeHTTP(rw, r.WithContext(ctx))
			return nil
		})
	}
}

func getVirtualWalletFromReq(r *http.Request) (*VirtualWallet, error) {
	vWallet, ok := r.Context().Value(virtualWalletKey{}).(*VirtualWallet)
	if !ok {
		return nil, fmt.Errorf("wallet not found")
	}
	return vWallet, nil
}

func checkBalance(next http.Handler) http.Handler {
	return wrapHandler(func(rw http.ResponseWriter, r *http.Request) error {
		vWallet, err := getVirtualWalletFromReq(r)
		if err != nil {
			return &ResponseError{
				StatusCode: http.StatusForbidden,
				WError:     err,
				Message:    "Transaction forbidden",
			}
		} else if vWallet.Balance <= 0 {
			return &ResponseError{
				StatusCode: http.StatusBadRequest,
				Message:    "Account balance is zero.",
			}
		}

		next.ServeHTTP(rw, r)
		return nil
	})
}
