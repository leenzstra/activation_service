package licensing

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/leenzstra/activation_service/internal/models"
	"github.com/leenzstra/activation_service/internal/utils"
)

type RegisterKeyBody struct {
	Key          string        `json:"key"`
	MaxUses      int           `json:"max_uses"`
	Contacts     string        `json:"contacts"`
	Expiration   time.Time     `json:"expiration"`
	SubjectsData []SubjectInfo `json:"subjects_data"`
}

type SubjectInfo struct {
	Sid     int      `json:"sid"`
	Classes []string `json:"classes"`
}

func getSubjectClasses(body *RegisterKeyBody, h *handler) ([]*models.SubjectClass, error) {
	s := []*models.SubjectClass{}
	for _, data := range body.SubjectsData {
		subjectClasses, err := h.DB.GetSubjectClasses(data.Sid, data.Classes)
		log.Println(subjectClasses, data, err)
		if err != nil {
			return nil, err
		}

		s = append(s, subjectClasses...)
	}
	return s, nil
}

func (h handler) RegisterKey(c *fiber.Ctx) error {

	log.Println("START")
	log.Println(string(c.Body()))

	payload := RegisterKeyBody{}
	if err := c.BodyParser(&payload); err != nil {
		return c.JSON(utils.WrapResponse(false, err.Error(), nil))
	}

	log.Println(payload)

	subjectClasses, err := getSubjectClasses(&payload, &h)
	if err != nil {
		return c.JSON(utils.WrapResponse(false, err.Error(), nil))
	}

	s, _ := json.MarshalIndent(subjectClasses, "", "\t")
	log.Print(string(s))

	err = h.DB.AddLicense(&models.License{
		Key:            payload.Key,
		MaxUses:        payload.MaxUses,
		Contacts:       payload.Contacts,
		SubjectClasses: subjectClasses,
		Expiration:     payload.Expiration,
	})
	if err != nil {
		return c.JSON(utils.WrapResponse(false, err.Error(), nil))
	}

	return c.JSON(utils.WrapResponse(true, "ok", nil))
}
