package main

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/nedpals/valentine-wall/backend/models"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/hook"
)

func internalError(data any) *apis.ApiError {
	return apis.NewApiError(
		http.StatusInternalServerError,
		http.StatusText(http.StatusInternalServerError),
		data,
	)
}

func setupRoutes(app *pocketbase.PocketBase) hook.Handler[*core.ServeEvent] {
	return func(e *core.ServeEvent) error {
		e.Router.Use(middleware.Recover())

		e.Router.Static("/renderer_assets", "renderer_assets")

		e.Router.GET("/terms-and-conditions", func(c echo.Context) error {
			// TODO
			return c.String(200, "TODO:")
		})

		e.Router.GET("/departments", func(c echo.Context) error {
			departments := []*models.CollegeDepartment{}
			err := models.DepartmentQuery(app.Dao()).All(&departments)
			if err != nil {
				return internalError(err)
			}

			return c.JSON(200, departments)
		})

		e.Router.GET("/gifts", func(c echo.Context) error {
			gifts := models.Gifts{}
			err := models.GiftQuery(app.Dao()).All(&gifts)
			if err != nil {
				return internalError(err)
			}

			return c.JSON(200, gifts)
		})

		return nil
	}
}
