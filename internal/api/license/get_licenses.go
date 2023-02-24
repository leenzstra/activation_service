package license

import (
	"github.com/gofiber/fiber/v2"
	"github.com/leenzstra/activation_service/internal/models"
	"github.com/leenzstra/activation_service/internal/utils"
)

type LicenseResponse struct {
	Key      string `json:"key"`
	MaxUses  int    `json:"max_uses"`
	Contacts string `json:"contacts"`
	Period   string `json:"period"`
	Uses     int    `json:"uses"`
}

func modelsToResponses(models []*models.License) []*LicenseResponse {
	responses := []*LicenseResponse{}

	for _, m := range models {
		responses = append(responses, &LicenseResponse{
			Key:      m.Key,
			MaxUses:  m.MaxUses,
			Contacts: m.Contacts,
			Period:   m.Period,
			Uses:     len(m.LicenseUses),
		})

	}
	return responses

}

func (h handler) GetLicenses(c *fiber.Ctx) error {
	licenses, err := h.DB.GetLicenses()
	if err != nil {
		return c.JSON(utils.WrapResponse(false, err.Error(), nil))
	}

	return c.JSON(utils.WrapResponse(true, "", modelsToResponses(licenses)))
}
