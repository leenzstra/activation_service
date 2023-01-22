package licensing

import (
	"encoding/json"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/hyperboloide/lk"
	"github.com/leenzstra/activation_service/internal/models"
	"github.com/leenzstra/activation_service/internal/utils"
)

type LicenseActivationBody struct {
	Key       string `json:"key"`
	MachineId string `json:"machine_id"`
}

type LicenseBody struct {
	CreationDate     string              `json:"creation_date"`
	Salt             string              `json:"salt"`
	MaxUses          int                 `json:"max_uses"`
	ExpirationDate   time.Time           `json:"exiration_date"`
	Contacts         string              `json:"contacts"`
	SelectedSubjects []*SubjectClassData `json:"selected_subjects"`
}

type SubjectClassData struct {
	SubjectName string `json:"subject_name"`
	Classes     string `json:"classes"`
	SubjectId   int    `json:"subject_id"`
	Alias       string `json:"alias"`
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

	if len(license.LicenseUses) >= license.MaxUses {
		return c.JSON(utils.WrapResponse(false, "лимит активаций ключа", nil))
	}

	licenseObject, err := lk.LicenseFromB32String(license.Key)
	if err != nil {
		return c.JSON(utils.WrapResponse(false, "license key decoding error", nil))
	}

	licenseData := LicenseBody{}
	if err := json.Unmarshal(licenseObject.Data, &licenseData); err != nil {
		return c.JSON(utils.WrapResponse(false, "license data decoding error", nil))
	}

	license_use := models.LicenseUse{LicenseID: license.ID, MachineId: payload.MachineId, ActivationDate: time.Now().UTC(), ExpirationDate: licenseData.ExpirationDate}
	if err := h.DB.AddLicenseUse(&license_use); err != nil {
		return c.JSON(utils.WrapResponse(false, "license activation error", nil))
	}

	return c.JSON(utils.WrapResponse(true, "ok", nil))
}
