package info

import (
	"github.com/gofiber/fiber/v2"
	"github.com/leenzstra/activation_service/internal/db"
)

type handler struct {
	DB *db.Database
}

func RegisterRoutes(app *fiber.App, db *db.Database) {
	h := &handler{
		DB: db,
	}

	routes := app.Group("/api/info")

	routes.Get("/subjects", h.GetSubjects)
}
