package models

import (
	"time"

	"gorm.io/gorm"
)

type License struct {
	gorm.Model
	Key            string `gorm:"unique"`
	MaxUses        int
	Contacts       string
	LicenseUses    []*LicenseUse
	SubjectClasses []*SubjectClass `gorm:"many2many:license_subject_classes;"`
	Period         string
}

type LicenseUse struct {
	gorm.Model
	LicenseID       uint
	Expiration      time.Time
	MachineInfoHash string
}

type Subject struct {
	gorm.Model
	Sid            int `gorm:"unique"`
	Name           string
	Alias          string `gorm:"unique"`
	SubjectClasses []*SubjectClass
}

type SubjectClass struct {
	gorm.Model
	SubjectID uint
	Class     string
}
