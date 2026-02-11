package main

import (
	"testing"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tests"
)

func TestOnAddUser(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	// Create user
	userCollection, _ := app.Dao().FindCollectionByNameOrId("users")
	user := models.NewRecord(userCollection)
	user.Set("email", "test@example.com")
	app.Dao().SaveRecord(user)
	
	e := &core.ModelEvent{
		Model: user,
	}
	
	err := onAddUser(app.Dao(), e)
	if err != nil {
		t.Errorf("onAddUser failed: %v", err)
	}
	
	// Verify wallet was created for user
	wallet, err := app.Dao().FindFirstRecordByData("virtual_wallets", "user", user.Id)
	if err != nil {
		t.Errorf("Wallet was not created for user: %v", err)
	}
	
	if wallet.GetFloat("balance") != 0.0 {
		t.Errorf("Expected initial balance 0, got %f", wallet.GetFloat("balance"))
	}
}

func TestOnAddUserDetails(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	// Create user
	userCollection, _ := app.Dao().FindCollectionByNameOrId("users")
	user := models.NewRecord(userCollection)
	user.Set("email", "test@example.com")
	app.Dao().SaveRecord(user)
	
	// Create user details
	detailsCollection, _ := app.Dao().FindCollectionByNameOrId("user_details")
	details := models.NewRecord(detailsCollection)
	details.Set("user", user.Id)
	details.Set("student_id", "202012345678")
	details.Set("email", "test@example.com")
	app.Dao().SaveRecord(details)
	
	e := &core.RecordCreateEvent{
		Record: details,
	}
	
	err := onAddUserDetails(app, e)
	if err != nil {
		t.Errorf("onAddUserDetails failed: %v", err)
	}
	
	// Verify user record was updated with details reference
	updatedUser, _ := app.Dao().FindRecordById("users", user.Id)
	if updatedUser.GetString("details") != details.Id {
		t.Error("User details reference was not set on user record")
	}
}

func TestOnRemoveUser(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	// Create user and details
	userCollection, _ := app.Dao().FindCollectionByNameOrId("users")
	user := models.NewRecord(userCollection)
	user.Set("email", "test@example.com")
	app.Dao().SaveRecord(user)
	
	detailsCollection, _ := app.Dao().FindCollectionByNameOrId("user_details")
	details := models.NewRecord(detailsCollection)
	details.Set("user", user.Id)
	details.Set("student_id", "202012345678")
	app.Dao().SaveRecord(details)
	
	e := &core.RecordDeleteEvent{
		Record: user,
	}
	
	err := onRemoveUser(app.Dao(), e)
	// Should attempt to delete user details
	_ = err // Error is acceptable if details not found
	
	// Verify user details were deleted (or attempted to be deleted)
	_, err = app.Dao().FindRecordById("user_details", details.Id)
	// Should get error as details should be deleted
	if err == nil {
		t.Error("User details were not deleted")
	}
}
