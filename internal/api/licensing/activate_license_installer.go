package licensing

import (
	"encoding/json"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/hyperboloide/lk"
	"github.com/leenzstra/activation_service/internal/models"
	"gorm.io/gorm"
)

func (h handler) ActivateLicenseInstaller(c *fiber.Ctx) error {
	payload := LicenseActivationBody{}
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(602).SendString("Parse data error")
	}

	license, err := h.DB.GetLicenseByKey(payload.Key)
	if err != nil {
		return c.Status(602).SendString("No that license key")
	}

	_, err = h.DB.GetLicenseUse(payload.Key, payload.MachineInfoHash)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		if len(license.LicenseUses) >= license.MaxUses {
			return c.Status(602).SendString("Превышен лимит активаций ключа")
		}

		license_use := &models.LicenseUse{
			LicenseID:       license.ID,
			MachineInfoHash: payload.MachineInfoHash,
		}

		if err := h.DB.AddLicenseUse(license_use); err != nil {
			return c.Status(602).SendString("Ошибка активации ключа")
		}
	} else if err != nil {
		return c.Status(602).SendString(err.Error())
	}

	licenseBody, _ := json.Marshal(LicenseBody{
		Key:             payload.Key,
		MachineInfoHash: payload.MachineInfoHash,
		Expiration:      license.Expiration,
	})

	licenseObj, err := lk.NewLicense(h.PrivateKey, licenseBody)
	if err != nil {
		return c.Status(602).SendString("Ошибка при создании лицензии")
	}

	licenseHash, err := licenseObj.ToHexString()
	if err != nil {
		return c.Status(602).SendString("Ошибка при преобразовании лицензии")
	}

	return c.Status(601).SendString(licenseHash)
}
