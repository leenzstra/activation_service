package license

import (
	"github.com/gofiber/fiber/v2"
	"github.com/leenzstra/activation_service/internal/models"
	"github.com/leenzstra/activation_service/internal/utils"
)

type AddKeyBody struct {
	Key          string        `json:"key"`
	MaxUses      int           `json:"max_uses"`
	Contacts     string        `json:"contacts"`
	Period       string        `json:"period"`
	SubjectsData []SubjectInfo `json:"subjects_data"`
}

type SubjectInfo struct {
	Sid     int      `json:"sid"`
	Classes []string `json:"classes"`
}

func getSubjectClasses(body *AddKeyBody, h *handler) ([]*models.SubjectClass, error) {
	s := []*models.SubjectClass{}

	for _, data := range body.SubjectsData {
		subjectClasses, err := h.DB.GetSubjectClasses(data.Sid, data.Classes)
		if err != nil {
			return nil, err
		}

		s = append(s, subjectClasses...)
	}
	return s, nil
}

func (h handler) RegisterKey(c *fiber.Ctx) error {
	payload := AddKeyBody{}
	if err := c.BodyParser(&payload); err != nil {
		return c.JSON(utils.WrapResponse(false, err.Error(), nil))
	}

	subjectClasses, err := getSubjectClasses(&payload, &h)
	if err != nil {
		return c.JSON(utils.WrapResponse(false, err.Error(), nil))
	}

	if err = h.DB.AddLicense(&models.License{
		Key:            payload.Key,
		MaxUses:        payload.MaxUses,
		Contacts:       payload.Contacts,
		SubjectClasses: subjectClasses,
		Period:         payload.Period,
	}); err != nil {
		return c.JSON(utils.WrapResponse(false, err.Error(), nil))
	}

	return c.JSON(utils.WrapResponse(true, "", nil))
}
