package main

import (
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
)

func createTransaction(dao *daos.Dao, wallet string, amount float64, description string) error {
	collection, err := dao.FindCollectionByNameOrId("virtual_transactions")
	if err != nil {
		return nil
	}

	record := models.NewRecord(collection)
	record.Set("wallet", wallet)
	record.Set("description", description)
	record.Set("amount", amount)
	return dao.SaveRecord(record)
}

func createTransactionFromUser(dao *daos.Dao, userId string, amount float64, description string) error {
	wallet, err := getWalletByUserId(dao, userId)
	if err != nil {
		return err
	}

	return createTransaction(dao, wallet.Id, amount, description)
}

func getWalletByUserId(dao *daos.Dao, userId string) (*models.Record, error) {
	return dao.FindFirstRecordByData("virtual_wallets", "user", userId)
}
