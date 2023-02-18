package licensing

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/hyperboloide/lk"
	"github.com/leenzstra/activation_service/internal/models"
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

func (h handler) ActivateLicense(c *fiber.Ctx) error {
	payload := LicenseActivationBody{}
	if err := c.BodyParser(&payload); err != nil {
		return c.JSON(utils.WrapResponse(false, err.Error(), nil))
	}

	license, err := h.DB.GetLicenseByKey(payload.Key)
	if err != nil {
		return c.JSON(utils.WrapResponse(false, err.Error(), nil))
	}

	_, err = h.DB.GetLicenseUse(payload.Key, payload.MachineInfoHash)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		if len(license.LicenseUses) >= license.MaxUses {
			return c.JSON(utils.WrapResponse(false, "превышен лимит активаций ключа", nil))
		}

		license_use := &models.LicenseUse{
			LicenseID:       license.ID,
			MachineInfoHash: payload.MachineInfoHash,
		}

		if err := h.DB.AddLicenseUse(license_use); err != nil {
			return c.JSON(utils.WrapResponse(false, "ошибка активации ключа", nil))
		}
	} else if err != nil {
		return c.JSON(utils.WrapResponse(false, err.Error(), nil))
	}

	licenseBody, _ := json.Marshal(LicenseBody{
		Key:             payload.Key,
		MachineInfoHash: payload.MachineInfoHash,
		Expiration:      license.Expiration,
	})

	licenseObj, err := lk.NewLicense(h.PrivateKey, licenseBody)
	if err != nil {
		return c.JSON(utils.WrapResponse(false, "ошибка при создании лицензии", nil))
	}

	licenseHash, err := licenseObj.ToHexString()
	if err != nil {
		return c.JSON(utils.WrapResponse(false, "ошибка при преобразовании лицензии", nil))
	}

	return c.JSON(utils.WrapResponse(true, "ok", licenseHash))
}
