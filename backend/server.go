package main

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	htmlTemplate "html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	vModels "github.com/nedpals/valentine-wall/backend/models"
	"github.com/patrickmn/go-cache"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/hook"
)

var zipFiles = sync.Map{}

func internalError(data any) *apis.ApiError {
	return apis.NewApiError(
		http.StatusInternalServerError,
		http.StatusText(http.StatusInternalServerError),
		data,
	)
}

var htmlTemplates = &htmlTemplate.Template{}

func setupRoutes(app *pocketbase.PocketBase) hook.Handler[*core.ServeEvent] {
	return func(e *core.ServeEvent) error {
		imageRenderer := &ImageRenderer{
			CacheStore: cache.New(time.Duration(10*time.Minute), time.Duration(5*time.Second)),
		}

		// chrome/browser-based image rendering specific code
		if len(chromeDevtoolsURL) != 0 {
			// launch chrome instance
			log.Printf("connecting chrome via: %s\n", chromeDevtoolsURL)
			remoteChromeCtx, remoteCtxCancel := chromedp.NewRemoteAllocator(context.Background(), chromeDevtoolsURL)
			defer remoteCtxCancel()

			chromeCtx, chromeCancel := chromedp.NewContext(remoteChromeCtx)
			defer chromeCancel()

			// load template
			log.Println("loading image template...")
			var err error
			if htmlTemplates, err = htmlTemplate.ParseGlob("./templates/html/*.html.tpl"); err != nil {
				log.Panicln(err)
			} else {
				log.Printf("%d html templates have been loaded\n", len(htmlTemplates.Templates()))
			}

			imageRenderer.ChromeCtx = chromeCtx
			imageRenderer.Template = htmlTemplates.Lookup("message_image.html.tpl")
		}

		tac, err := getTermsAndConditions()
		if err != nil {
			return err
		}

		e.Router.Use(middleware.Recover())

		e.Router.Static("/renderer_assets", "renderer_assets")

		e.Router.GET("/terms-and-conditions", func(c echo.Context) error {
			_, err := c.Response().Write(tac)
			return err
		})

		e.Router.GET("/departments", func(c echo.Context) error {
			departments := []*vModels.CollegeDepartment{}
			err := vModels.DepartmentQuery(app.Dao()).All(&departments)
			if err != nil {
				return internalError(err)
			}

			return c.JSON(200, departments)
		})

		e.Router.GET("/gifts", func(c echo.Context) error {
			gifts := vModels.Gifts{}
			err := vModels.GiftQuery(app.Dao()).All(&gifts)
			if err != nil {
				return internalError(err)
			}

			return c.JSON(200, gifts)
		})

		e.Router.GET("/messages/{messageId}/image", func(c echo.Context) error {
			id := c.PathParam("messageId")
			message, err := app.Dao().FindRecordById("messages", id)
			if err != nil {
				return apis.NewNotFoundError("Message not found", err)
			}

			query := c.QueryParams()
			if query.Has("template_image") {
				c.Response().Header().Set("Content-Type", "text/html")
				if err := imageRenderer.Template.Execute(c.Response(), RendererContext{
					RawMessage: message,
					BackendURL: baseUrl,
				}); err != nil {
					return err
				}
				return nil
			}

			buf, err := imageRenderer.Render(imageTypeTwitter, message)
			if err != nil {
				return err
			}

			c.Response().Header().Set("Content-Type", "image/png")
			c.Response().Write(buf)
			return nil
		})

		e.Router.GET("/user_messages/archive", func(c echo.Context) error {
			authRecord := c.Get(apis.ContextAuthRecordKey).(*models.Record)
			authDetails := authRecord.Expand()["details"].(*models.Record)

			// get recipient id
			recipientId := authDetails.GetString("student_id")

			c.Response().Header().Set("Content-Type", "text/event-stream")
			c.Response().Header().Set("Cache-Control", "no-cache")
			c.Response().Header().Set("Connection", "keep-alive")

			connCloseChan := make(chan struct{})
			defer close(connCloseChan)

			errChan := make(chan error)
			defer close(errChan)

			rw := c.Response().Writer
			encodeDataSSE(rw, map[string]any{"status": "starting"})

			go func() {
				// TODO:
				// stats, err := getRecipientStatsBySID(messagesSearchIndex, recipientId)
				// if err != nil {
				// 	errChan <- err
				// 	return
				// }

				// encodeDataSSE(rw, map[string]any{
				// 	"status": "set_file_count",
				// 	"data": map[string]int{
				// 		"count": stats.GiftMessagesCount + stats.MessagesCount,
				// 	},
				// })

				_, loaded := zipFiles.LoadAndDelete(recipientId)
				if !loaded {
					messages, err := app.Dao().FindRecordsByExpr("messages", dbx.HashExp{"recipient": recipientId})
					if err != nil {
						errChan <- err
						return
					}

					zipArchive := &bytes.Buffer{}
					zipWriter := zip.NewWriter(zipArchive)

					for _, msg := range messages {
						encodeDataSSE(rw, map[string]any{
							"status": "processing",
							"data":   map[string]int{"len": 1},
						})

						buf, err := imageRenderer.Render(imageTypeTwitter, msg)
						if err != nil {
							errChan <- err
							return
						}

						filePath := filepath.Join("messages", fmt.Sprintf("messages_%s_%s.png", msg.GetString("recipient"), msg.Id))
						fileWriter, err := zipWriter.Create(filePath)
						if err != nil {
							errChan <- err
							return
						}

						fileWriter.Write(buf)
					}

					// TODO: generate summary
					passivePrintError(zipWriter.Close())
					zipFiles.Store(recipientId, zipArchive.Bytes())
				}

				// TODO: add archive stats
				// if _, err := db.NamedExec(
				// 	"INSERT INTO archived_stats (email, total) VALUES (:email, :total) ON CONFLICT (email) DO NOTHING",
				// 	&ArchiveStats{gotUser.Email, stats.GiftMessagesCount + stats.MessagesCount},
				// ); err != nil {
				// 	passivePrintError(err)
				// }

				encodeDataSSE(rw, map[string]any{
					"status": "done",
					"data": map[string]string{
						"endpoint": "/user_messages/download_archive",
					},
				})

				<-connCloseChan
			}()

			for {
				select {
				case err := <-errChan:
					encodeDataSSE(rw, map[string]any{
						"status": "error",
						"data":   map[string]string{"message": err.Error()},
					})
					<-connCloseChan
				case <-connCloseChan:
					return nil
				case <-c.Request().Context().Done():
					return nil
				}
			}
		}, apis.RequireRecordAuth("users"))

		e.Router.GET("/user_messages/download_archive", func(c echo.Context) error {
			authRecord := c.Get(apis.ContextAuthRecordKey).(*models.Record)
			authDetails := authRecord.Expand()["details"].(*models.Record)

			// get email
			recipientId := authDetails.GetString("student_id")
			gotZip, loaded := zipFiles.LoadAndDelete(recipientId)
			if !loaded {
				return apis.NewNotFoundError("Not found", nil)
			}

			zipName := fmt.Sprintf("archive_%s.zip", recipientId)
			c.Response().Header().Set("Content-Type", "application/zip")
			c.Response().Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", zipName))
			_, err := c.Response().Write(gotZip.([]byte))
			return err
		}, apis.RequireRecordAuth("users"))

		return nil
	}
}
