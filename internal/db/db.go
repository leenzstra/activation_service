package db

import (
	"log"
	"os"
	"time"

	"github.com/leenzstra/activation_service/internal/config"
	"github.com/leenzstra/activation_service/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	Gorm *gorm.DB
	Config *config.Config
}

func New(config config.Config) *Database {
	db, err := gorm.Open(sqlite.Open(config.DBPath), &gorm.Config{Logger: logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
		SlowThreshold:             time.Second, // Slow SQL threshold
		LogLevel:                  logger.Info, // Log level
		IgnoreRecordNotFoundError: false,       // Ignore ErrRecordNotFound error for logger
		Colorful:                  true,        // Disable color
	})})

	if err != nil {
		log.Fatalln(err)
	}

	db.AutoMigrate(&models.License{}, &models.LicenseUse{})

	return &Database{Gorm: db, Config: &config}
}

func (d Database) AddLicense(license *models.License) error {
	return d.Gorm.Create(&license).Error
}

func (d Database) AddLicenseUse(use *models.LicenseUse) error {
	return d.Gorm.Create(use).Error
}

func (d Database) GetLicenseByKey(key string) (*models.License, error) {
	license := models.License{}
	if err := d.Gorm.Preload("LicenseUses").Where("key = ?", key).First(&license).Error; err != nil {
		return nil, err
	}

	return &license, nil
}

func (d Database) GetLicenses() ([]*models.License, error) {
	licenses := []*models.License{}
	if err := d.Gorm.Preload("LicenseUses").Find(&licenses).Error; err != nil {
		return nil, err
	}

	return licenses, nil
}
