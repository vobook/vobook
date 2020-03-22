package tests

import (
	"testing"
	"time"

	"vobook/cmd/server/requests"
	"vobook/cmd/server/responses"
	"vobook/database/factories"
	"vobook/database/models"
	contactpropertytype "vobook/enum/contact_property_type"
	. "vobook/tests/apitest"
	"vobook/tests/assert"
	"vobook/utils"

	fake "github.com/brianvoe/gofakeit"
)

func TestUpdateContactProperty(t *testing.T) {
	u := Login(t)
	contact, err := factories.CreateContact(models.Contact{UserID: u.ID})
	assert.NotError(t, err)
	prop, err := factories.CreateContactProperty(models.ContactProperty{ContactID: contact.ID})
	assert.NotError(t, err)

	name := fake.Name()
	value := fake.Name()
	req := requests.UpdateContactProperty{
		Name:  &name,
		Value: &value,
	}

	var resp responses.Success
	TestUpdate(t, "contact-properties/"+prop.ID, req, &resp)

	assert.DatabaseHas(t, "contact_properties", utils.M{
		"id":    prop.ID,
		"type":  prop.Type,
		"name":  name,
		"value": value,
	})
}

func TestTrashContactProperties(t *testing.T) {
	u := Login(t)
	contact, err := factories.CreateContact(models.Contact{UserID: u.ID})
	assert.NotError(t, err)

	prop1, err := factories.CreateContactProperty(models.ContactProperty{ContactID: contact.ID})
	assert.NotError(t, err)
	prop2, err := factories.CreateContactProperty(models.ContactProperty{ContactID: contact.ID})
	assert.NotError(t, err)
	prop3, err := factories.CreateContactProperty(models.ContactProperty{ContactID: contact.ID})
	assert.NotError(t, err)

	req := requests.IDs{
		prop1.ID,
		prop2.ID,
	}
	var resp responses.Success
	TestUpdate(t, "trash-contact-properties", req, &resp)

	assert.DatabaseHasDeleted(t, "contact_properties", prop1.ID, prop2.ID)
	assert.DatabaseHas(t, "contact_properties", utils.M{
		"id":         prop3.ID,
		"deleted_at": nil,
	})
}

func TestRestoreContactProperties(t *testing.T) {
	u := Login(t)
	contact, err := factories.CreateContact(models.Contact{UserID: u.ID})
	assert.NotError(t, err)

	deletedAt := time.Now()
	prop1, err := factories.CreateContactProperty(models.ContactProperty{
		ContactID: contact.ID,
		DeletedAt: &deletedAt,
	})
	assert.NotError(t, err)
	prop2, err := factories.CreateContactProperty(models.ContactProperty{
		ContactID: contact.ID,
		DeletedAt: &deletedAt,
	})
	assert.NotError(t, err)
	prop3, err := factories.CreateContactProperty(models.ContactProperty{
		ContactID: contact.ID,
		DeletedAt: &deletedAt,
	})
	assert.NotError(t, err)

	req := requests.IDs{
		prop1.ID,
		prop2.ID,
	}
	var resp responses.Success
	TestUpdate(t, "restore-contact-properties", req, &resp)

	assert.DatabaseHas(t, "contact_properties", utils.M{
		"id":         prop1.ID,
		"deleted_at": nil,
	})
	assert.DatabaseHas(t, "contact_properties", utils.M{
		"id":         prop2.ID,
		"deleted_at": nil,
	})
	assert.DatabaseHas(t, "contact_properties", utils.M{
		"id":         prop3.ID,
		"deleted_at": assert.NotNil,
	})
}

func TestDeleteContactProperties(t *testing.T) {
	u := Login(t)
	contact, err := factories.CreateContact(models.Contact{UserID: u.ID})
	assert.NotError(t, err)

	deletedAt := time.Now()
	prop1, err := factories.CreateContactProperty(models.ContactProperty{
		ContactID: contact.ID,
		DeletedAt: &deletedAt,
	})
	assert.NotError(t, err)
	prop2, err := factories.CreateContactProperty(models.ContactProperty{
		ContactID: contact.ID,
	})
	assert.NotError(t, err)
	prop3, err := factories.CreateContactProperty(models.ContactProperty{
		ContactID: contact.ID,
	})
	assert.NotError(t, err)

	req := requests.IDs{
		prop1.ID,
		prop2.ID,
	}
	var resp responses.Success
	TestUpdate(t, "delete-contact-properties", req, &resp)

	assert.DatabaseMissing(t, "contact_properties", utils.M{
		"id": prop1.ID,
	})
	assert.DatabaseMissing(t, "contact_properties", utils.M{
		"id": prop2.ID,
	})
	assert.DatabaseHas(t, "contact_properties", utils.M{
		"id": prop3.ID,
	})
}

func TestReorderContactProperties(t *testing.T) {
	u := Login(t)
	contact, err := factories.CreateContact(models.Contact{UserID: u.ID})
	assert.NotError(t, err)

	prop1, err := factories.CreateContactProperty(models.ContactProperty{
		ContactID: contact.ID,
		Order:     10,
	})
	assert.NotError(t, err)
	prop2, err := factories.CreateContactProperty(models.ContactProperty{
		ContactID: contact.ID,
		Order:     4,
	})
	assert.NotError(t, err)
	prop3, err := factories.CreateContactProperty(models.ContactProperty{
		ContactID: contact.ID,
		Order:     12,
	})
	assert.NotError(t, err)

	req := requests.IDs{
		prop1.ID,
		prop2.ID,
		prop3.ID,
	}
	var resp responses.Success
	TestUpdate(t, "reorder-contact-properties", req, &resp)

	assert.DatabaseHas(t, "contact_properties", utils.M{
		"id":    prop1.ID,
		"order": 0,
	})
	assert.DatabaseHas(t, "contact_properties", utils.M{
		"id":    prop2.ID,
		"order": 1,
	})
	assert.DatabaseHas(t, "contact_properties", utils.M{
		"id":    prop3.ID,
		"order": 2,
	})
}

func TestGetContactPropertyTypes(t *testing.T) {
	Login(t)

	var resp []contactpropertytype.TypeModel
	TestGet(t, "contact-property-types", &resp)
	assert.Equals(t, len(contactpropertytype.All), len(resp))
}
