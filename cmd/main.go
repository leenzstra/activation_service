package main

import (
	"flag"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/hyperboloide/lk"
	"github.com/leenzstra/activation_service/internal/api/content"
	"github.com/leenzstra/activation_service/internal/api/info"
	"github.com/leenzstra/activation_service/internal/api/license"
	"github.com/leenzstra/activation_service/internal/config"
	"github.com/leenzstra/activation_service/internal/db"
	"github.com/leenzstra/activation_service/internal/keypair"
)

const pub = ""
const priv = ""

func main() {
    isProd :=flag.Bool("prod", false, "")
	flag.Parse()

	cfg, err := config.LoadConfig(*isProd)
	if err != nil {
		log.Fatalln("Failed at cfg", err)
	}

	subjectsConfig, err := config.LoadSubjectsConfig()
	if err != nil {
		log.Fatalln("Failed at subjects cfg", err)
	}

	pubKey, err := lk.PublicKeyFromB32String(pub)
    if err != nil {
		log.Fatalln("Failed at pub key creation", err)
	}

	privateKey, err := lk.PrivateKeyFromB32String(priv)
	if err != nil {
		log.Fatalln("Failed at private key creation", err)
	}

    app := fiber.New()
    app.Use(logger.New())
	app.Use(recover.New())

    db := db.New(cfg)
	if db.InitSubjectsInfo(&subjectsConfig) != nil {
		log.Fatalln("Failed at init subjects info", err)
	}

	license.RegisterRoutes(app, db, keypair.New(pubKey, privateKey))
	info.RegisterRoutes(app, db)
	content.RegisterRoutes(app, db)

    app.Get("/api/status", func(ctx *fiber.Ctx) error {
        return ctx.SendString("ok")
    })

    app.Listen(":3000")
}