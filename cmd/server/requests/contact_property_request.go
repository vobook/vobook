package requests

import (
	"github.com/vovainside/vobook/database/models"
	"github.com/vovainside/vobook/enum/contact_property"
)

type CreateContactProperty struct {
	Type  contactproperty.Type `json:"type"`
	Name  string               `json:"name"`
	Value string               `json:"value"`
	Order int                  `json:"order"`
}

func (r *CreateContactProperty) Validate() (err error) {
	err = r.Type.Validate()
	return
}

func (r *CreateContactProperty) ToModel() *models.ContactProperty {
	return &models.ContactProperty{
		Type:  r.Type,
		Name:  r.Name,
		Value: r.Value,
		Order: r.Order,
	}
}

type UpdateContactProperty struct {
	Name  *string `json:"name"`
	Value *string `json:"value"`
}

func (r *UpdateContactProperty) Validate() (err error) {
	return
}

func (r *UpdateContactProperty) ToModel(m *models.ContactProperty) {
	if r.Name != nil {
		m.Name = *r.Name
	}
	if r.Value != nil {
		m.Value = *r.Value
	}
}
