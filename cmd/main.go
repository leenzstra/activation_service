package main

import (
	"flag"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/hyperboloide/lk"
	"github.com/leenzstra/activation_service/internal/api/licensing"
	"github.com/leenzstra/activation_service/internal/config"
	"github.com/leenzstra/activation_service/internal/db"
)

const pub = "AS73O4LLJYHWD7V2WCWBZBDAYFQCEWM6TJQ2WEUJKVARKD4EIG5XNGOBQ5QMSMB343AXW7LNX4P4MJMU6XA27574WFQVUY6RY4J2LZMT2MOBGJZUUIGOIBYTCN2WO36FTIKJT3JFSXOIQ3MRH6XK2VSKX3FA===="

func main() {

    isProd :=flag.Bool("prod", false, "")
	flag.Parse()

	c, err := config.LoadConfig(*isProd)

	if err != nil {
		log.Fatalln("Failed at config", err)
	}

    app := fiber.New()
    app.Use(logger.New())
	app.Use(recover.New())

    db := db.New(c)

    pubKey, err := lk.PublicKeyFromB32String(pub)
    if err != nil {
		log.Fatalln("Failed at pub key creation", err)
	}

	licensing.RegisterRoutes(app, db, pubKey)

    app.Get("/api/status", func(c *fiber.Ctx) error {
        return c.SendString("ok")
    })

    app.Listen(":3000")
}