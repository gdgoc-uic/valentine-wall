package main

import (
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
)

func onAddWallet(dao *daos.Dao, e *core.ModelEvent) error {
	// add initial balance
	return createTransaction(dao, e.Model.GetId(), 1000, "Initial balance")
}

func onAddWalletTransaction(dao *daos.Dao, e *core.ModelEvent) error {
	// add transaction amount to wallet
	transaction, err := dao.FindRecordById(e.Model.TableName(), e.Model.GetId())
	if err != nil {
		return err
	}

	record, err := dao.FindRecordById("virtual_wallets", transaction.GetString("wallet"))
	if err != nil {
		return err
	}

	record.Set("balance", record.GetFloat("balance")+transaction.GetFloat("amount"))
	return dao.SaveRecord(record)
}
