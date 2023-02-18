package access

import (
	"log"

	"github.com/gofiber/fiber/v2"
)


type Config struct {
	Filter       func(c *fiber.Ctx) bool // Required
	Unauthorized fiber.Handler           // middleware specfic
	Validator    func(*fiber.Ctx) (bool, error)

	SuccessHandler fiber.Handler

	ErrorHandler fiber.ErrorHandler
}


func New(config ...Config) fiber.Handler {
	var cfg Config
	if len(config) > 0 {
		cfg = config[0]
	}

	if cfg.SuccessHandler == nil {
		cfg.SuccessHandler = func(c *fiber.Ctx) error {
			log.Println("ACCESS GRANTED")
			return c.Next()
		}
	}
	if cfg.ErrorHandler == nil {
		cfg.ErrorHandler = func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusUnauthorized).SendString("Access denied")
		}
	}
	// Return middleware handler
	return func(c *fiber.Ctx) error {
		// Filter request to skip middleware
		if cfg.Filter != nil && cfg.Filter(c) {
			return c.Next()
		}

		valid, err := cfg.Validator(c)

		if err == nil && valid {
			return cfg.SuccessHandler(c)
		}
		return cfg.ErrorHandler(c, err)
	}
}
