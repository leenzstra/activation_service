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
const priv = "FD7YCAYBAEFXA22DN5XHIYLJNZSXEAP7QIAACAQBANIHKYQBBIAACAKEAH7YIAAAAAFP7AYFAEBP7BQAAAAP7GP7QIAWCBF7W5YWWTQPMH7LVMFMDSCGBQLAEJMZ5GTBVMJISVKBCUHYIQN3O2M4DB3AZEYDXZWBPN6W3PY7YYSZJ5OBV737ZMLBLJR5DRYTUXSZHUY4CMTTJIQM4QDRGE3VM5X4LGQUTHWSLFO4RBWZCP5OVVLEVPWKAEYQEHAIQEJG4NZWVZC6VHR35DJOKE22KKQYDNGQKC54NMOMZIHJG2QZUBPR26EGWF2Z56ZX33LIL2C46IAA===="

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

	licensing.RegisterRoutes(app, db, pubKey, privateKey)

    app.Get("/api/status", func(ctx *fiber.Ctx) error {
        return ctx.SendString("ok")
    })

    app.Listen(":3000")
}