package license

import (
	"errors"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/keyauth/v2"
	"github.com/leenzstra/activation_service/internal/db"
	"github.com/leenzstra/activation_service/internal/keypair"
	"github.com/leenzstra/activation_service/internal/middlewares/access"
)

var (
	ErrLicenseCreation = errors.New("ошибка при создании лицензии")
)

type handler struct {
	DB      *db.Database
	KeyPair *keypair.KeyPair
}

func RegisterRoutes(app *fiber.App, db *db.Database, keyPair *keypair.KeyPair) {
	h := &handler{
		DB:         db,
		KeyPair: keyPair,
	}

	var authMiddleware = keyauth.New(keyauth.Config{
		Validator: func(c *fiber.Ctx, key string) (bool, error) {
			if key == h.DB.Config.ApiKey {
				return true, nil
			}
			return false, fiber.ErrUnauthorized
		},
	})

	var keyInfoAccessMiddleware = access.New(access.Config{
		Validator: func(c *fiber.Ctx) (bool, error) {
			payload := LicenseActivationBody{}
			if err := c.BodyParser(&payload); err != nil {
				return false, fiber.ErrForbidden
			}
			_, err := h.DB.GetLicenseUse(payload.Key, payload.MachineInfoHash)
			auth := strings.Split(c.Get("Authorization"), " ")
			token := ""
			if len(auth) == 2 {
				token = auth[1]
			}
			log.Println(payload.Key, payload.MachineInfoHash, token, err, h.DB.Config.ApiKey == token)
			if h.DB.Config.ApiKey == token || err == nil {
				return true, nil
			} else {
				return false, fiber.ErrUnauthorized
			}
		},
	})

	routes := app.Group("/api/license")

	routes.Post("/activate", h.ActivateLicense)
	routes.Post("/register", authMiddleware, h.RegisterKey)
	routes.Post("/verify", h.VerifyLicense)
	routes.Post("/info", keyInfoAccessMiddleware, h.GetLicenseInfo)
	routes.Get("/all", authMiddleware, h.GetLicenses)
}
