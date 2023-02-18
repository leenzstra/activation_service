package licensing

import (
	"github.com/gofiber/fiber/v2"
	"github.com/leenzstra/activation_service/internal/collections"
	"github.com/leenzstra/activation_service/internal/models"
	"github.com/leenzstra/activation_service/internal/utils"
)

type SubjectResponse struct {
	Sid     int      `json:"sid"`
	Name    string   `json:"name"`
	Alias   string   `json:"alias"`
	Classes []string `json:"classes"`
}

func modelsToSubjectResponses(subjects []*models.Subject) []*SubjectResponse {

	getClasses := func(classes []*models.SubjectClass) []string {
		return collections.Map(classes, func(c *models.SubjectClass) string {
			return c.Class
		})
	}

	return collections.Map(subjects, func(o *models.Subject) *SubjectResponse {
		return &SubjectResponse{
			Sid:     o.Sid,
			Name:    o.Name,
			Alias:   o.Alias,
			Classes: getClasses(o.SubjectClasses),
		}
	})

}

func (h handler) GetSubjects(c *fiber.Ctx) error {
	subjects, err := h.DB.GetSubjects()
	if err != nil {
		return c.JSON(utils.WrapResponse(false, err.Error(), nil))
	}

	return c.JSON(utils.WrapResponse(true, "ok", modelsToSubjectResponses(subjects)))
}
