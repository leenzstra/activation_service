package licensing

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/hyperboloide/lk"
	"github.com/leenzstra/activation_service/internal/models"
	"github.com/leenzstra/activation_service/internal/utils"
)

type AddLicenseBody struct {
	Key     string `json:"key"`
}

func (h handler) AddLicense(c *fiber.Ctx) error {
	payload := AddLicenseBody{}
	if err := c.BodyParser(&payload); err != nil {
		return c.JSON(utils.WrapResponse(false, err.Error(), nil))
	}

	licenseObject, err := lk.LicenseFromB32String(payload.Key)
	if err != nil {
		return c.JSON(utils.WrapResponse(false, "license key decoding error", nil))
	}

	licenseData := LicenseBody{}
	if err := json.Unmarshal(licenseObject.Data, &licenseData); err != nil {
		return c.JSON(utils.WrapResponse(false, "license data decoding error", nil))
	}

	err = h.DB.AddLicense(&models.License{Key: payload.Key, MaxUses: licenseData.MaxUses, Contacts: licenseData.Contacts})
	if err != nil {
		return c.JSON(utils.WrapResponse(false, err.Error(), nil))
	}

	return c.JSON(utils.WrapResponse(true, "ok", nil))
}
