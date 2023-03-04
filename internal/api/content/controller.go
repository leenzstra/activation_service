package content

import (
	"github.com/gofiber/fiber/v2"
	"github.com/leenzstra/activation_service/internal/db"
	"github.com/gofiber/keyauth/v2"
)

type handler struct {
	DB *db.Database
}

func RegisterRoutes(app *fiber.App, db *db.Database) {
	h := &handler{
		DB: db,
	}

	var licenseAuth = keyauth.New(keyauth.Config{
		KeyLookup: "header:License",
		Validator: func(c *fiber.Ctx, key string) (bool, error) {
			if key == h.DB.Config.ApiKey {
				return true, nil
			}
			return false, fiber.ErrUnauthorized
		},
	})

	routes := app.Group("/api/content")
	files := routes.Group("/files", licenseAuth)

	files.Static("/", "./files")
	// routes.Get("/subjects", h.GetSubjects)
}
