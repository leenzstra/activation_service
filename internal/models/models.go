package models

import (
	"time"

	"gorm.io/gorm"
)

type License struct {
	gorm.Model
	Key         string `gorm:"unique"`
	MaxUses     int
	Contacts    string
	LicenseUses []LicenseUse
}

type LicenseUse struct {
	gorm.Model
	LicenseID      uint
	MachineId      string
	ActivationDate time.Time
	ExpirationDate time.Time
}
