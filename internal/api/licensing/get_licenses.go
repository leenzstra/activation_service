package licensing

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/leenzstra/activation_service/internal/models"
	"github.com/leenzstra/activation_service/internal/utils"
)

type LicenseResponse struct {
	Key         string                `json:"key"`
	MaxUses     int                   `json:"max_uses"`
	Contacts    string                `json:"contacts"`
	LicenseUses []*LicenseUseResponse `json:"uses"`
}

type LicenseUseResponse struct {
	MachineId string    `json:"machine_id"`
	Activation    time.Time `json:"activation"`
	Expiration    time.Time `json:"expiration"`
}

func modelsToResponses(models []*models.License) []*LicenseResponse {
	responses := []*LicenseResponse{}

	for _, m := range models {
		uses := []*LicenseUseResponse{}
		for _, u := range m.LicenseUses {
			uses = append(uses, &LicenseUseResponse{MachineId: u.MachineId, Activation: u.ActivationDate, Expiration: u.ExpirationDate})
		}
		responses = append(responses, &LicenseResponse{Key: m.Key, MaxUses: m.MaxUses, LicenseUses: uses, Contacts: m.Contacts})

	}
	return responses

}

func (h handler) GetLicenses(c *fiber.Ctx) error {
	licenses, err := h.DB.GetLicenses()
	if err != nil {
		return c.JSON(utils.WrapResponse(false, err.Error(), nil))
	}

	return c.JSON(utils.WrapResponse(true, "ok", modelsToResponses(licenses)))
}
