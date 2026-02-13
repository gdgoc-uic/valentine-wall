package main

import (
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
)

func onAddUser(dao *daos.Dao, e *core.ModelEvent) error {
	// create virtual wallet
	collection, err := dao.FindCollectionByNameOrId("virtual_wallets")
	if err != nil {
		return nil
	}

	record := models.NewRecord(collection)
	record.Set("user", e.Model.GetId())
	record.Set("balance", 0)

	return dao.SaveRecord(record)
}

func onAddUserDetails(app core.App, e *core.RecordCreateEvent) error {
	dao := app.Dao()
	user, err := dao.FindRecordById("users", e.Record.GetString("user"))
	if err != nil {
		return err
	}

	user.Set("details", e.Record.Id)
	passivePrintError(dao.SaveRecord(user))

	// TODO: build message count
	// messageCount, giftMessagesCount := 0, 0

	// if res, err := dao.FindRecordsByExpr("messages", dbx.HashExp{"recipient": e.Record.GetString("student_id")}); err == nil {
	// 	for _, row := range res {

	// 	}
	// }

	// TODO: add message count
	email := e.Record.Email()
	if msg, err := emailTemplates.welcome.With(map[string]any{
		"Email": email,
	}).Message(app.Settings().Meta, email); err == nil {
		passivePrintError(app.NewMailClient().Send(msg))
	}

	return nil
}

func onUserVerified(app core.App, e *core.RecordConfirmVerificationEvent) error {
	// TODO: add message count
	msg, err := emailTemplates.welcome.Message(app.Settings().Meta, e.Record.Email())
	if err != nil {
		// TODO: add error
		// return err
	}

	return app.NewMailClient().Send(msg)
}

func onRemoveUser(dao *daos.Dao, e *core.RecordDeleteEvent) error {
	record, err := dao.FindFirstRecordByData("user_details", "user", e.Record.Id)
	if err != nil {
		// TODO: add error
		// return nil
	}

	return dao.DeleteRecord(record)
}
