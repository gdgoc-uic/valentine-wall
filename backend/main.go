package main

import (
	"context"
	"log"

	"github.com/chromedp/chromedp"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"

	_ "github.com/nedpals/valentine-wall/backend/migrations"
)

func main() {
	app := pocketbase.New()
	migratecmd.MustRegister(app, app.RootCmd, &migratecmd.Options{
		Automigrate: true,
	})

	// chrome/browser-based image rendering specific code
	if len(chromeDevtoolsURL) != 0 {
		// launch chrome instance
		log.Printf("connecting chrome via: %s\n", chromeDevtoolsURL)
		remoteChromeCtx, remoteCtxCancel := chromedp.NewRemoteAllocator(context.Background(), chromeDevtoolsURL)
		defer remoteCtxCancel()

		chromeCtx, chromeCancel := chromedp.NewContext(remoteChromeCtx)
		defer chromeCancel()

		imageRenderer.ChromeCtx = chromeCtx
	}

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
		case "user_details":
			return onAddUserDetails(app.Dao(), e)
		case "messages":
			return onAddMessage(app.Dao(), e)
		case "message_replies":
			return onAddMessageReply(app.Dao(), e)
		}

		return nil
	})

	app.OnModelAfterCreate().Add(func(e *core.ModelEvent) error {
		switch e.Model.TableName() {
		case "users":
			return onAddUser(app.Dao(), e)
		case "virtual_wallets":
			return onAddWallet(app.Dao(), e)
		case "virtual_transactions":
			return onAddWalletTransaction(app.Dao(), e)
		}

		return nil
	})

	app.OnRecordAfterDeleteRequest().Add(func(e *core.RecordDeleteEvent) error {
		switch e.Record.Collection().Name {
		case "users":
			return onRemoveUser(app.Dao(), e)
		case "messages":
			return onRemoveMessage(app.Dao(), e)
		case "message_replies":
			// NOTE: temp added
			return onRemoveMessageReply(app.Dao(), e)
		}

		return nil
	})

	app.OnBeforeServe().Add(setupRoutes(app))

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
