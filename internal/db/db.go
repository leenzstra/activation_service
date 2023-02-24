package db

import (
	"log"
	"os"
	"time"

	"github.com/leenzstra/activation_service/internal/collections"
	"github.com/leenzstra/activation_service/internal/config"
	"github.com/leenzstra/activation_service/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	Gorm   *gorm.DB
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

	if db.AutoMigrate(&models.License{}, &models.LicenseUse{}, &models.Subject{}, &models.SubjectClass{}) != nil {
		log.Fatalln(err)
	}

	return &Database{Gorm: db, Config: &config}
}

func (d Database) InitSubjectsInfo(sc *config.SubjectsConfig) error {
	subjectsCount := int64(0)
	if err := d.Gorm.Model(&models.Subject{}).Count(&subjectsCount).Error; err != nil {
		return err
	}

	if subjectsCount == 0 {
		for _, subject := range sc.Subjects {
			subjModel := models.Subject{Sid: subject.Sid, Name: subject.Name, Alias: subject.Alias}
			d.Gorm.Create(&subjModel)
			subjClasses := collections.Map(subject.Classes, func(s string) models.SubjectClass { return models.SubjectClass{SubjectID: subjModel.ID, Class: s} })
			d.Gorm.Create(&subjClasses)
		}
		return nil
	} else {
		return nil
	}
}

func (d Database) AddLicense(license *models.License) error {
	return d.Gorm.Create(&license).Error
}

func (d Database) AddLicenseUse(use *models.LicenseUse) error {
	return d.Gorm.Create(use).Error
}

func (d Database) GetLicenseByKey(key string) (*models.License, error) {
	license := models.License{}
	if err := d.Gorm.Preload("LicenseUses").Preload("SubjectClasses").Where("key = ?", key).First(&license).Error; err != nil {
		return nil, err
	}

	return &license, nil
}

func (d Database) GetSubjects() ([]*models.Subject, error) {
	subjects := []*models.Subject{}
	if err := d.Gorm.Preload("SubjectClasses").Find(&subjects).Error; err != nil {
		return nil, err
	}
	return subjects, nil
}

func (d Database) GetSubjectsBySid(sids ...int) ([]*models.Subject, error) {
	subjects := []*models.Subject{}
	if err := d.Gorm.Preload("SubjectClasses").Where("sid in (?)", sids).Find(&subjects).Error; err != nil {
		return nil, err
	}
	return subjects, nil
}

func (d Database) GetSubjectsById(ids ...uint) ([]*models.Subject, error) {
	subjects := []*models.Subject{}
	if err := d.Gorm.Preload("SubjectClasses").Where("id in (?)", ids).Find(&subjects).Error; err != nil {
		return nil, err
	}
	return subjects, nil
}

func (d Database) GetLicenseUse(key, machineInfoHash string) (*models.LicenseUse, error) {
	license := &models.License{}
	if err := d.Gorm.
		Preload("LicenseUses", "machine_info_hash = ?", machineInfoHash).
		Where("key = ?", key).
		Find(&license).Error; err != nil {
		return nil, err
	}
	if len(license.LicenseUses) == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return license.LicenseUses[0], nil
}

func (d Database) GetSubjectClasses(sid int, classes []string) ([]*models.SubjectClass, error) {
	subjectClasses := []*models.SubjectClass{}

	if err := d.Gorm.
		Raw("SELECT subject_classes.id, subject_classes.created_at,subject_classes.updated_at, subject_classes.deleted_at, subjects.id as subject_id, class FROM 'subjects' JOIN 'subject_classes' ON 'subjects'.'id' = 'subject_classes'.'subject_id' WHERE sid = ? AND class IN (?);", sid, classes).
		Find(&subjectClasses).Error; err != nil {
		return nil, err
	}
	return subjectClasses, nil
}

func (d Database) GetLicenses() ([]*models.License, error) {
	licenses := []*models.License{}
	if err := d.Gorm.
		Preload("LicenseUses").
		Find(&licenses).Error; err != nil {
		return nil, err
	}

	return licenses, nil
}
