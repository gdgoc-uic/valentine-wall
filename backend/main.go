package main

import (
	"log"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
)

func main() {
	app := pocketbase.New()
	migratecmd.MustRegister(app, app.RootCmd, &migratecmd.Options{
		Automigrate: true,
	})

	app.OnRecordAfterConfirmVerificationRequest().Add(func(e *core.RecordConfirmVerificationEvent) error {
		return onUserVerified(app, e)
	})

	app.OnRecordBeforeCreateRequest().Add(func(e *core.RecordCreateEvent) error {
		switch e.Record.Collection().Name {
		case "messages":
			return onBeforeAddMessage(app.Dao(), e)
		}

		return nil
	})

	app.OnRecordAfterCreateRequest().Add(func(e *core.RecordCreateEvent) error {
		switch e.Record.Collection().Name {
		case "users":
			return onAddUser(app.Dao(), e)
		case "messages":
			return onAddMessage(app.Dao(), e)
		case "message_replies":
			return onAddMessageReply(app.Dao(), e)
		}

		return nil
	})

	app.OnModelAfterCreate().Add(func(e *core.ModelEvent) error {
		switch e.Model.TableName() {
		case "virtual_wallets":
			return onAddWallet(app.Dao(), e)
		case "virtual_transactions":
			return onAddWalletTransaction(app.Dao(), e)
		}

		return nil
	})

	app.OnRecordAfterDeleteRequest().Add(func(e *core.RecordDeleteEvent) error {
		switch e.Record.Collection().Name {
		case "messages":
			return onRemoveMessage(app.Dao(), e)
		}

		return nil
	})

	app.OnBeforeServe().Add(setupRoutes(app))

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
