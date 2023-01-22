package licensing

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/keyauth/v2"
	"github.com/hyperboloide/lk"
	"github.com/leenzstra/activation_service/internal/db"
)

type handler struct {
	DB *db.Database
	PubKey *lk.PublicKey
}

func RegisterRoutes(app *fiber.App, db *db.Database, pubKey *lk.PublicKey) {
	h := &handler{
		DB: db,
		PubKey: pubKey,
	}

	var authMiddleware = keyauth.New(keyauth.Config{
		Validator:  func(c *fiber.Ctx, key string) (bool, error) {
			if key == h.DB.Config.ApiKey {
				return true, nil
			}
			return false, fiber.ErrUnauthorized
		},
	  })

	routes := app.Group("/api/license")
	routes.Post("/activate",  h.ActivateLicense)
	routes.Post("/add", authMiddleware, h.AddLicense)
	routes.Get("/all", authMiddleware, h.GetLicenses)
}
