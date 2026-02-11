package main

import (
	"testing"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tests"
)

func TestCheckSufficientFunds_Success(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	// Create user and wallet
	userCollection, _ := app.Dao().FindCollectionByNameOrId("users")
	user := models.NewRecord(userCollection)
	app.Dao().SaveRecord(user)
	
	walletCollection, _ := app.Dao().FindCollectionByNameOrId("virtual_wallets")
	wallet := models.NewRecord(walletCollection)
	wallet.Set("user", user.Id)
	wallet.Set("balance", 1000.0)
	app.Dao().SaveRecord(wallet)
	
	err := checkSufficientFunds(app.Dao(), user.Id, 500.0)
	if err != nil {
		t.Errorf("Expected no error for sufficient funds, got: %v", err)
	}
}

func TestCheckSufficientFunds_Insufficient(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	// Create user and wallet with low balance
	userCollection, _ := app.Dao().FindCollectionByNameOrId("users")
	user := models.NewRecord(userCollection)
	app.Dao().SaveRecord(user)
	
	walletCollection, _ := app.Dao().FindCollectionByNameOrId("virtual_wallets")
	wallet := models.NewRecord(walletCollection)
	wallet.Set("user", user.Id)
	wallet.Set("balance", 100.0)
	app.Dao().SaveRecord(wallet)
	
	err := checkSufficientFunds(app.Dao(), user.Id, 500.0)
	if err == nil {
		t.Error("Expected insufficient funds error, got nil")
	}
}

func TestCreateTransaction(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	// Create wallet
	walletCollection, _ := app.Dao().FindCollectionByNameOrId("virtual_wallets")
	wallet := models.NewRecord(walletCollection)
	wallet.Set("balance", 1000.0)
	app.Dao().SaveRecord(wallet)
	
	// Create transaction
	amount := -150.0
	description := "Test transaction"
	
	err := createTransaction(app.Dao(), wallet.Id, amount, description)
	if err != nil {
		t.Errorf("createTransaction failed: %v", err)
	}
	
	// Verify transaction was created
	transactions, err := app.Dao().FindRecordsByExpr("virtual_transactions", 
		nil,
	)
	if err != nil || len(transactions) == 0 {
		t.Error("Transaction was not created")
	}
}

func TestOnAddWallet(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	// Create wallet
	walletCollection, _ := app.Dao().FindCollectionByNameOrId("virtual_wallets")
	wallet := models.NewRecord(walletCollection)
	wallet.Set("balance", 0.0)
	app.Dao().SaveRecord(wallet)
	
	e := &core.ModelEvent{
		Model: wallet,
	}
	
	err := onAddWallet(app.Dao(), e)
	if err != nil {
		t.Errorf("onAddWallet failed: %v", err)
	}
	
	// Verify initial balance transaction was created
	transactions, err := app.Dao().FindRecordsByExpr("virtual_transactions", nil)
	if err != nil || len(transactions) == 0 {
		t.Error("Initial balance transaction was not created")
	}
	
	// Check that transaction description contains "Initial balance"
	found := false
	for _, tx := range transactions {
		if tx.GetString("wallet") == wallet.Id && tx.GetFloat("amount") == 1000.0 {
			found = true
			break
		}
	}
	
	if !found {
		t.Error("Initial balance of 1000 coins was not added")
	}
}

func TestOnAddWalletTransaction(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	// Create wallet
	walletCollection, _ := app.Dao().FindCollectionByNameOrId("virtual_wallets")
	wallet := models.NewRecord(walletCollection)
	wallet.Set("balance", 1000.0)
	app.Dao().SaveRecord(wallet)
	
	// Create transaction
	txCollection, _ := app.Dao().FindCollectionByNameOrId("virtual_transactions")
	transaction := models.NewRecord(txCollection)
	transaction.Set("wallet", wallet.Id)
	transaction.Set("amount", -150.0)
	transaction.Set("description", "Test deduction")
	app.Dao().SaveRecord(transaction)
	
	e := &core.ModelEvent{
		Model: transaction,
	}
	
	err := onAddWalletTransaction(app.Dao(), e)
	if err != nil {
		t.Errorf("onAddWalletTransaction failed: %v", err)
	}
	
	// Verify wallet balance was updated
	updatedWallet, _ := app.Dao().FindRecordById("virtual_wallets", wallet.Id)
	expectedBalance := 850.0 // 1000 - 150
	
	if updatedWallet.GetFloat("balance") != expectedBalance {
		t.Errorf("Expected balance %f, got %f", expectedBalance, updatedWallet.GetFloat("balance"))
	}
}

func TestGetWalletByUserId(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	// Create user and wallet
	userCollection, _ := app.Dao().FindCollectionByNameOrId("users")
	user := models.NewRecord(userCollection)
	app.Dao().SaveRecord(user)
	
	walletCollection, _ := app.Dao().FindCollectionByNameOrId("virtual_wallets")
	wallet := models.NewRecord(walletCollection)
	wallet.Set("user", user.Id)
	wallet.Set("balance", 1000.0)
	app.Dao().SaveRecord(wallet)
	
	// Get wallet by user ID
	foundWallet, err := getWalletByUserId(app.Dao(), user.Id)
	if err != nil {
		t.Errorf("getWalletByUserId failed: %v", err)
	}
	
	if foundWallet.Id != wallet.Id {
		t.Errorf("Expected wallet ID %s, got %s", wallet.Id, foundWallet.Id)
	}
}

func TestCreateTransactionFromUser(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	// Create user and wallet
	userCollection, _ := app.Dao().FindCollectionByNameOrId("users")
	user := models.NewRecord(userCollection)
	app.Dao().SaveRecord(user)
	
	walletCollection, _ := app.Dao().FindCollectionByNameOrId("virtual_wallets")
	wallet := models.NewRecord(walletCollection)
	wallet.Set("user", user.Id)
	wallet.Set("balance", 1000.0)
	app.Dao().SaveRecord(wallet)
	
	// Create transaction from user
	err := createTransactionFromUser(app.Dao(), user.Id, 50.0, "Test reward")
	if err != nil {
		t.Errorf("createTransactionFromUser failed: %v", err)
	}
	
	// Check transaction was created
	transactions, _ := app.Dao().FindRecordsByExpr("virtual_transactions", nil)
	found := false
	for _, tx := range transactions {
		if tx.GetString("wallet") == wallet.Id && tx.GetFloat("amount") == 50.0 {
			found = true
			break
		}
	}
	
	if !found {
		t.Error("Transaction from user was not created")
	}
}
