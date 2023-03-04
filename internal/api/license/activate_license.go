package license

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/hyperboloide/lk"
	"github.com/leenzstra/activation_service/internal/keypair"
	"github.com/leenzstra/activation_service/internal/models"
	"github.com/leenzstra/activation_service/internal/responses"
	"github.com/leenzstra/activation_service/internal/utils"
	"gorm.io/gorm"
)

type LicenseActivationBody struct {
	Key             string `json:"key"`
	MachineInfoHash string `json:"machine_info_hash"`
}

type LicenseBody struct {
	Key             string    `json:"key"`
	MachineInfoHash string    `json:"machine_info_hash"`
	Expiration      time.Time `json:"expiration"`
}

func generateLicenseHash(keypair *keypair.KeyPair, body *LicenseBody) (string, error) {

	licenseBody, err := json.Marshal(body)
	if err != nil {
		return "", ErrLicenseCreation
	}

	licenseObj, err := lk.NewLicense(keypair.Private, licenseBody)
	if err != nil {
		return "", ErrLicenseCreation
	}

	licenseHash, err := licenseObj.ToHexString()
	if err != nil {
		return "", ErrLicenseCreation
	}

	return licenseHash, nil
}

func (h handler) ActivateLicense(c *fiber.Ctx) error {
	payload := LicenseActivationBody{}
	if err := c.BodyParser(&payload); err != nil {
		return c.JSON(utils.WrapResponse(false, err.Error(), nil))
	}

	license, err := h.DB.GetLicenseByKey(payload.Key)
	if err != nil {
		return c.JSON(utils.WrapResponse(false, err.Error(), nil))
	}

	license_use, err := h.DB.GetLicenseUse(payload.Key, payload.MachineInfoHash)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		if len(license.LicenseUses) >= license.MaxUses {
			return c.JSON(utils.WrapResponse(false, "превышен лимит активаций ключа", nil))
		}

		period := &responses.LicenseDuration{}
		if err := period.FromString(license.Period); err != nil {
			return c.JSON(utils.WrapResponse(false, err.Error(), nil))
		}

		license_use := &models.LicenseUse{
			LicenseID:       license.ID,
			MachineInfoHash: payload.MachineInfoHash,
			Expiration:      time.Now().UTC().AddDate(period.Years, period.Months, period.Days),
		}

		if err := h.DB.AddLicenseUse(license_use); err != nil {
			return c.JSON(utils.WrapResponse(false, "ошибка активации ключа", nil))
		}
	} else if err != nil {
		return c.JSON(utils.WrapResponse(false, err.Error(), nil))
	} 

	licenseBody := LicenseBody{
		Key:             payload.Key,
		MachineInfoHash: payload.MachineInfoHash,
		Expiration:      license_use.Expiration,
	}

	license_use, err = h.DB.GetLicenseUse(payload.Key, payload.MachineInfoHash)
	if err != nil {
		return c.JSON(utils.WrapResponse(false, err.Error(), nil))
	}
	
	hash, err := generateLicenseHash(h.KeyPair, &licenseBody)
	if err != nil {
		return c.JSON(utils.WrapResponse(false, err.Error(), nil))
	}
	
	fmt.Println(licenseBody, hash)
	return c.JSON(utils.WrapResponse(true, "", hash))
}
