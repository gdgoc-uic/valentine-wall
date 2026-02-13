package main

import (
	"testing"

	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
	"github.com/pocketbase/pocketbase/tests"
)

// newTestApp creates a test app with the project's custom collections.
func newTestApp(t *testing.T) *tests.TestApp {
	t.Helper()

	app, err := tests.NewTestApp()
	if err != nil {
		t.Fatalf("Failed to create test app: %v", err)
	}

	dao := app.Dao()

	// Create "messages" collection
	messages := &models.Collection{}
	messages.Name = "messages"
	messages.Type = models.CollectionTypeBase
	messages.Schema = schema.NewSchema(
		&schema.SchemaField{Name: "content", Type: schema.FieldTypeText},
		&schema.SchemaField{Name: "recipient", Type: schema.FieldTypeText},
		&schema.SchemaField{Name: "user", Type: schema.FieldTypeText},
		&schema.SchemaField{Name: "replies_count", Type: schema.FieldTypeNumber},
		&schema.SchemaField{Name: "gifts", Type: schema.FieldTypeJson},
	)
	if err := dao.SaveCollection(messages); err != nil {
		app.Cleanup()
		t.Fatalf("Failed to create messages collection: %v", err)
	}

	// Create "user_details" collection
	userDetails := &models.Collection{}
	userDetails.Name = "user_details"
	userDetails.Type = models.CollectionTypeBase
	userDetails.Schema = schema.NewSchema(
		&schema.SchemaField{Name: "user", Type: schema.FieldTypeText},
		&schema.SchemaField{Name: "student_id", Type: schema.FieldTypeText},
		&schema.SchemaField{Name: "email", Type: schema.FieldTypeText},
		&schema.SchemaField{Name: "last_active", Type: schema.FieldTypeDate},
	)
	if err := dao.SaveCollection(userDetails); err != nil {
		app.Cleanup()
		t.Fatalf("Failed to create user_details collection: %v", err)
	}

	// Create "virtual_wallets" collection
	wallets := &models.Collection{}
	wallets.Name = "virtual_wallets"
	wallets.Type = models.CollectionTypeBase
	wallets.Schema = schema.NewSchema(
		&schema.SchemaField{Name: "user", Type: schema.FieldTypeText},
		&schema.SchemaField{Name: "balance", Type: schema.FieldTypeNumber},
	)
	if err := dao.SaveCollection(wallets); err != nil {
		app.Cleanup()
		t.Fatalf("Failed to create virtual_wallets collection: %v", err)
	}

	// Create "virtual_transactions" collection
	transactions := &models.Collection{}
	transactions.Name = "virtual_transactions"
	transactions.Type = models.CollectionTypeBase
	transactions.Schema = schema.NewSchema(
		&schema.SchemaField{Name: "wallet", Type: schema.FieldTypeText},
		&schema.SchemaField{Name: "amount", Type: schema.FieldTypeNumber},
		&schema.SchemaField{Name: "description", Type: schema.FieldTypeText},
	)
	if err := dao.SaveCollection(transactions); err != nil {
		app.Cleanup()
		t.Fatalf("Failed to create virtual_transactions collection: %v", err)
	}

	// Create "message_replies" collection
	replies := &models.Collection{}
	replies.Name = "message_replies"
	replies.Type = models.CollectionTypeBase
	replies.Schema = schema.NewSchema(
		&schema.SchemaField{Name: "content", Type: schema.FieldTypeText},
		&schema.SchemaField{Name: "sender", Type: schema.FieldTypeRelation, Options: &schema.RelationOptions{
			CollectionId: userDetails.Id,
			MaxSelect:    ptrInt(1),
		}},
		&schema.SchemaField{Name: "message", Type: schema.FieldTypeRelation, Options: &schema.RelationOptions{
			CollectionId: messages.Id,
			MaxSelect:    ptrInt(1),
		}},
		&schema.SchemaField{Name: "liked", Type: schema.FieldTypeBool},
	)
	if err := dao.SaveCollection(replies); err != nil {
		app.Cleanup()
		t.Fatalf("Failed to create message_replies collection: %v", err)
	}

	return app
}

func ptrInt(i int) *int {
	return &i
}

