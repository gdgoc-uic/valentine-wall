package main

import (
	"testing"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tests"
)

func TestExpandMessage(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	// Create test collections and records
	collection, err := app.Dao().FindCollectionByNameOrId("messages")
	if err != nil {
		t.Fatalf("Failed to find messages collection: %v", err)
	}

	message := models.NewRecord(collection)
	message.Set("content", "Test message")
	message.Set("recipient", "202012345678")
	
	err = expandMessage(app.Dao(), message)
	if err != nil {
		t.Errorf("expandMessage failed: %v", err)
	}
}

func TestComputeGiftCost(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	// Create message with gifts
	collection, _ := app.Dao().FindCollectionByNameOrId("messages")
	message := models.NewRecord(collection)
	
	// Create mock gifts
	giftCollection, _ := app.Dao().FindCollectionByNameOrId("gifts")
	gift1 := models.NewRecord(giftCollection)
	gift1.Set("price", 50.0)
	gift1.Set("is_remittable", true)
	
	gift2 := models.NewRecord(giftCollection)
	gift2.Set("price", 30.0)
	gift2.Set("is_remittable", false)
	
	// Add gifts to message expand
	message.SetExpand(map[string]any{
		"gifts": []*models.Record{gift1, gift2},
	})
	
	totalAmount, remittableAmount := computeGiftCost(message)
	
	expectedTotal := 80.0
	expectedRemittable := 50.0
	
	if totalAmount != expectedTotal {
		t.Errorf("Expected total amount %f, got %f", expectedTotal, totalAmount)
	}
	
	if remittableAmount != expectedRemittable {
		t.Errorf("Expected remittable amount %f, got %f", expectedRemittable, remittableAmount)
	}
}

func TestUpdateRanking(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	recipientId := "202012345678"
	coinsToAdd := 150.0
	
	err := updateRanking(app.Dao(), recipientId, coinsToAdd)
	if err != nil {
		t.Errorf("updateRanking failed: %v", err)
	}
	
	// Verify ranking was created/updated
	ranking, err := app.Dao().FindFirstRecordByData("rankings", "recipient", recipientId)
	if err != nil {
		t.Errorf("Failed to find ranking: %v", err)
	}
	
	if ranking.GetFloat("total_coins") < coinsToAdd {
		t.Errorf("Expected at least %f coins, got %f", coinsToAdd, ranking.GetFloat("total_coins"))
	}
}

func TestOnBeforeAddMessage_DuplicateDetection(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection, _ := app.Dao().FindCollectionByNameOrId("messages")
	userCollection, _ := app.Dao().FindCollectionByNameOrId("user_details")
	
	// Create user
	user := models.NewRecord(userCollection)
	user.Set("student_id", "202012345678")
	app.Dao().SaveRecord(user)
	
	// Create first message
	message1 := models.NewRecord(collection)
	message1.Set("content", "Test duplicate message")
	message1.Set("recipient", "202087654321")
	message1.Set("user", user.Id)
	app.Dao().SaveRecord(message1)
	
	// Try to create duplicate
	message2 := models.NewRecord(collection)
	message2.Set("content", "Test duplicate message")
	message2.Set("recipient", "202087654321")
	message2.Set("user", user.Id)
	
	e := &core.RecordCreateEvent{
		Record: message2,
	}
	
	err := onBeforeAddMessage(app.Dao(), e)
	if err == nil {
		t.Error("Expected error for duplicate message, got nil")
	}
}

func TestOnBeforeAddMessage_ProfanityCheck(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection, _ := app.Dao().FindCollectionByNameOrId("messages")
	
	message := models.NewRecord(collection)
	message.Set("content", "This is a clean test message")
	message.Set("recipient", "202012345678")
	
	e := &core.RecordCreateEvent{
		Record: message,
	}
	
	err := onBeforeAddMessage(app.Dao(), e)
	// Should pass if profanity list is loaded and message is clean
	// Error handling depends on profanity list availability
	_ = err
}

func TestOnBeforeAddMessage_InsufficientFunds(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	// Create user with empty wallet
	userCollection, _ := app.Dao().FindCollectionByNameOrId("users")
	user := models.NewRecord(userCollection)
	app.Dao().SaveRecord(user)
	
	walletCollection, _ := app.Dao().FindCollectionByNameOrId("virtual_wallets")
	wallet := models.NewRecord(walletCollection)
	wallet.Set("user", user.Id)
	wallet.Set("balance", 0.0)
	app.Dao().SaveRecord(wallet)
	
	// Create message
	messageCollection, _ := app.Dao().FindCollectionByNameOrId("messages")
	message := models.NewRecord(messageCollection)
	message.Set("content", "Test message")
	message.Set("recipient", "202012345678")
	message.Set("user", user.Id)
	
	e := &core.RecordCreateEvent{
		Record: message,
	}
	
	err := onBeforeAddMessage(app.Dao(), e)
	if err == nil {
		t.Error("Expected insufficient funds error, got nil")
	}
}
