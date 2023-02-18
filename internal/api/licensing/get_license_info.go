package licensing

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/leenzstra/activation_service/internal/collections"
	"github.com/leenzstra/activation_service/internal/models"
	"github.com/leenzstra/activation_service/internal/utils"
)

type LicenseInfoResponse struct {
	Key          string                 `json:"key"`
	MaxUses      int                    `json:"max_uses"`
	Contacts     string                 `json:"contacts"`
	Expiration   time.Time              `json:"expiration"`
	LicenseUses  []*LicenseUseResponse  `json:"uses"`
	SubjectsInfo []*SubjectInfoResponse `json:"subjects_info"`
}

type SubjectInfoResponse struct {
	Sid     uint     `json:"sid"`
	Name    string   `json:"name"`
	Alias   string   `json:"alias"`
	Classes []string `json:"classes"`
}

type LicenseUseResponse struct {
	MachineInfoHash string `json:"machine_info_hash"`
}

func modelsToLicenseInfoResponses(license *models.License, selectedSubjects []*models.Subject) *LicenseInfoResponse {

	uses := []*LicenseUseResponse{}
	for _, use := range license.LicenseUses {
		uses = append(uses, &LicenseUseResponse{MachineInfoHash: use.MachineInfoHash})
	}

	subjects := []*SubjectInfoResponse{}
	for _, subject := range selectedSubjects {
		classes := []string{}
		for _, class := range license.SubjectClasses {
			if class.SubjectID == subject.ID {
				classes = append(classes, class.Class)
			}
		}
		subjects = append(subjects, &SubjectInfoResponse{
			Sid:     uint(subject.Sid),
			Name:    subject.Name,
			Classes: classes,
			Alias:   subject.Alias,
		})
	}

	response := &LicenseInfoResponse{
		Key:          license.Key,
		MaxUses:      license.MaxUses,
		Contacts:     license.Contacts,
		Expiration:   license.Expiration,
		LicenseUses:  uses,
		SubjectsInfo: subjects,
	}

	return response

}

func (h handler) GetLicenseInfo(c *fiber.Ctx) error {
	payload := LicenseActivationBody{}
	if err := c.BodyParser(&payload); err != nil {
		return c.JSON(utils.WrapResponse(false, err.Error(), nil))
	}

	license, err := h.DB.GetLicenseByKey(payload.Key)
	if err != nil {
		return c.JSON(utils.WrapResponse(false, err.Error(), nil))
	}

	ids := collections.Map(license.SubjectClasses, func(sc *models.SubjectClass) uint { return sc.SubjectID })

	subjects, err := h.DB.GetSubjectsById(ids...)
	if err != nil {
		return c.JSON(utils.WrapResponse(false, err.Error(), nil))
	}

	return c.JSON(utils.WrapResponse(true, "ok", modelsToLicenseInfoResponses(license, subjects)))
}
