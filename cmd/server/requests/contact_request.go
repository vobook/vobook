package requests

import (
	"time"

	"github.com/vovainside/vobook/database/models"
)

type CreateContact struct {
	Name       string                  `json:"name"`
	FirstName  string                  `json:"first_name"`
	LastName   string                  `json:"last_name"`
	MiddleName string                  `json:"middle_name"`
	Birthday   time.Time               `json:"birthday"`
	Properties []CreateContactProperty `json:"properties"`
}

func (r *CreateContact) Validate() (err error) {
	for _, v := range r.Properties {
		err = v.Validate()
		if err != nil {
			return
		}
	}
	return
}

func (r *CreateContact) ToModel() *models.Contact {
	m := &models.Contact{
		Name:       r.Name,
		FirstName:  r.FirstName,
		LastName:   r.LastName,
		MiddleName: r.MiddleName,
		Birthday:   r.Birthday,
	}

	m.Properties = make([]*models.ContactProperty, len(r.Properties))
	for i, v := range r.Properties {
		m.Properties[i] = v.ToModel()
	}

	return m
}
