package license

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/hyperboloide/lk"
	"github.com/leenzstra/activation_service/internal/utils"
)

type LicenseVerificationBody struct {
	LicenseHash     string `json:"license_hash"`
	MachineInfoHash string `json:"machine_info_hash"`
}

type VerifiedBody struct {
	Status     int `json:"status"`
}

func (h handler) VerifyLicense(c *fiber.Ctx) error {
	payload := LicenseVerificationBody{}
	if err := c.BodyParser(&payload); err != nil {
		return c.JSON(utils.WrapResponse(false, err.Error(), nil))
	}

	licenseObj, err := lk.LicenseFromHexString(payload.LicenseHash)
	if err != nil {
		return c.JSON(utils.WrapResponse(false, "ошибка при создании лицензии", nil))
	}

	if ok, err := licenseObj.Verify(h.KeyPair.Public); err != nil {
		return c.JSON(utils.WrapResponse(false, "ошибка верификации лицензии", nil))
	} else if !ok {
		return c.JSON(utils.WrapResponse(false, "лицензия недействительна", nil))
	}

	result := LicenseBody{}

	if err := json.Unmarshal(licenseObj.Data, &result); err != nil {
		return c.JSON(utils.WrapResponse(false, err.Error(), nil))
	}

	_, err = h.DB.GetLicenseByKey(result.Key)
	if err != nil {
		return c.JSON(utils.WrapResponse(false, err.Error(), nil))
	}

	if (result.MachineInfoHash != payload.MachineInfoHash){
		return c.JSON(utils.WrapResponse(false, "лицензия активирована на другой компьютер", nil))
	}

	return c.JSON(utils.WrapResponse(true, "", VerifiedBody{Status: 1}))
}
