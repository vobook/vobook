package requests

import (
	"vobook/database/models"
	contactpropertytype "vobook/enum/contact_property_type"
)

type CreateContactProperty struct {
	Type  contactpropertytype.Type `json:"type"`
	Name  string                   `json:"name"`
	Value string                   `json:"value"`
}

func (r *CreateContactProperty) Validate() (err error) {
	err = r.Type.Validate()
	return
}

func (r *CreateContactProperty) ToModel() models.ContactProperty {
	return models.ContactProperty{
		Type:  r.Type,
		Name:  r.Name,
		Value: r.Value,
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
