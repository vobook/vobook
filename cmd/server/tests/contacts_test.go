package tests

import (
	"fmt"
	"testing"
	"time"
	"vobook/domain/file"
	"vobook/services/fs"

	"vobook/cmd/server/requests"
	"vobook/cmd/server/responses"
	"vobook/database/factories"
	"vobook/database/models"
	contactpropertytype "vobook/enum/contact_property_type"
	. "vobook/tests/apitest"
	"vobook/tests/assert"
	. "vobook/tests/fake"
	"vobook/utils"

	fake "github.com/brianvoe/gofakeit"
)

func TestCreateContact(t *testing.T) {
	req := requests.CreateContact{
		Name:      fake.Name(),
		FirstName: fake.FirstName(),
		LastName:  fake.LastName(),
		DOBYear:   fake.Year(),
		DOBMonth:  time.Month(fake.Number(1, 12)),
		DOBDay:    fake.Day(),
		Properties: []requests.CreateContactProperty{
			{
				Type:  contactpropertytype.Email,
				Value: fake.Email(),
			},
			{
				Type:  contactpropertytype.Phone,
				Value: fake.PhoneFormatted(),
			},
			{
				Type:  contactpropertytype.Phone,
				Value: fake.PhoneFormatted(),
				Name:  "Only SMS",
			},
			{
				Type:  contactpropertytype.Link,
				Value: fake.URL(),
				Name:  "limit for me",
			},
			{
				Type:  contactpropertytype.Other,
				Value: fake.Sentence(3),
				Name:  "hello",
			},
		},
	}

	var resp models.Contact
	TestCreate(t, "contacts", req, &resp)

	assert.Equals(t, AuthUser.ID, resp.UserID)
	assert.Equals(t, req.FirstName, resp.FirstName)
	assert.Equals(t, req.LastName, resp.LastName)
	assert.Equals(t, req.DOBYear, resp.DOBYear)
	assert.Equals(t, req.DOBMonth, resp.DOBMonth)
	assert.Equals(t, req.DOBDay, resp.DOBDay)
	assert.Equals(t, len(req.Properties), len(resp.Props))

	for i, v := range req.Properties {
		assert.Equals(t, resp.ID, resp.Props[i].ContactID)
		assert.Equals(t, v.Name, resp.Props[i].Name)
		assert.Equals(t, v.Value, resp.Props[i].Value)
		assert.Equals(t, v.Type, resp.Props[i].Type)
	}

	assert.DatabaseHas(t, "contacts", utils.M{
		"id":         resp.ID,
		"user_id":    AuthUser.ID,
		"first_name": req.FirstName,
		"last_name":  req.LastName,
		"dob_year":   req.DOBYear,
		"dob_month":  req.DOBMonth,
		"dob_day":    req.DOBDay,
	})

	for _, v := range req.Properties {
		assert.DatabaseHas(t, "contact_properties", utils.M{
			"contact_id": resp.ID,
			"name":       v.Name,
			"value":      v.Value,
			"type":       v.Type,
		})
	}
}

func TestUpdateContact(t *testing.T) {
	Login(t)
	elem, err := factories.CreateContact(models.Contact{UserID: AuthUser.ID})
	assert.NotError(t, err)

	name := fake.Name()
	firstName := fake.FirstName()
	lastName := fake.LastName()
	middleName := fake.LastName()
	dobYear := fake.Year()
	dobMonth := time.Month(fake.Number(1, 12))
	dobDay := fake.Day()
	req := requests.UpdateContact{
		Name:       &name,
		FirstName:  &firstName,
		LastName:   &lastName,
		MiddleName: &middleName,
		DOBYear:    &dobYear,
		DOBMonth:   &dobMonth,
		DOBDay:     &dobDay,
		Properties: []requests.CreateContactProperty{
			{
				Type:  contactpropertytype.Email,
				Value: fake.Email(),
			},
			{
				Type:  contactpropertytype.Phone,
				Value: fake.PhoneFormatted(),
			},
		},
	}

	var resp responses.Success
	TestUpdate(t, "contacts/"+elem.ID, req, &resp)

	assert.DatabaseHas(t, "contacts", utils.M{
		"id":          elem.ID,
		"name":        name,
		"first_name":  firstName,
		"last_name":   lastName,
		"middle_name": middleName,
		"dob_year":    dobYear,
		"dob_month":   dobMonth,
		"dob_day":     dobDay,
	})

	assert.DatabaseHas(t, "contact_properties", utils.M{
		"contact_id": elem.ID,
		"type":       contactpropertytype.Email,
		"value":      req.Properties[0].Value,
	})
	assert.DatabaseHas(t, "contact_properties", utils.M{
		"contact_id": elem.ID,
		"type":       contactpropertytype.Phone,
		"value":      req.Properties[1].Value,
	})
}

func TestTrashContacts(t *testing.T) {
	Login(t)
	elem1, err := factories.CreateContact(models.Contact{UserID: AuthUser.ID})
	assert.NotError(t, err)
	elem2, err := factories.CreateContact(models.Contact{UserID: AuthUser.ID})
	assert.NotError(t, err)
	elem3, err := factories.CreateContact(models.Contact{UserID: AuthUser.ID})
	assert.NotError(t, err)

	req := requests.IDs{
		elem1.ID,
		elem2.ID,
	}
	var resp responses.Success
	TestUpdate(t, "trash-contacts", req, &resp)

	assert.DatabaseHasDeleted(t, "contacts", elem1.ID, elem2.ID)
	assert.DatabaseHas(t, "contacts", utils.M{
		"id":         elem3.ID,
		"deleted_at": nil,
	})
}

func TestRestoreContacts(t *testing.T) {
	Login(t)
	deletedAt := time.Now()
	elem1, err := factories.CreateContact(models.Contact{UserID: AuthUser.ID, DeletedAt: &deletedAt})
	assert.NotError(t, err)
	elem2, err := factories.CreateContact(models.Contact{UserID: AuthUser.ID, DeletedAt: &deletedAt})
	assert.NotError(t, err)
	elem3, err := factories.CreateContact(models.Contact{UserID: AuthUser.ID, DeletedAt: &deletedAt})
	assert.NotError(t, err)

	req := requests.IDs{
		elem1.ID,
		elem2.ID,
	}
	var resp responses.Success
	TestUpdate(t, "restore-contacts", req, &resp)

	assert.DatabaseHasDeleted(t, "contacts", elem3.ID)
	assert.DatabaseHas(t, "contacts", utils.M{
		"id":         elem1.ID,
		"deleted_at": nil,
	})
	assert.DatabaseHas(t, "contacts", utils.M{
		"id":         elem2.ID,
		"deleted_at": nil,
	})
}

func TestDeleteContacts(t *testing.T) {
	Login(t)
	elem1, err := factories.CreateContact(models.Contact{UserID: AuthUser.ID})
	assert.NotError(t, err)
	elem2, err := factories.CreateContact(models.Contact{UserID: AuthUser.ID})
	assert.NotError(t, err)
	elem3, err := factories.CreateContact(models.Contact{UserID: AuthUser.ID})
	assert.NotError(t, err)

	req := requests.IDs{
		elem1.ID,
		elem2.ID,
	}
	var resp responses.Success
	TestUpdate(t, "delete-contacts", req, &resp)

	assert.DatabaseHas(t, "contacts", utils.M{
		"id": elem3.ID,
	})
	assert.DatabaseMissing(t, "contacts", utils.M{
		"id": elem1.ID,
	})
	assert.DatabaseMissing(t, "contacts", utils.M{
		"id": elem2.ID,
	})
}

func TestGetContact(t *testing.T) {
	Login(t)
	elem, err := factories.CreateContact(models.Contact{UserID: AuthUser.ID})
	assert.NotError(t, err)

	props := make([]models.ContactProperty, 3)
	for i := range props {
		prop, err := factories.CreateContactProperty(models.ContactProperty{
			ContactID: elem.ID,
			Name:      fmt.Sprintf("Prop %d", i+1),
			Order:     i + 1,
		})
		assert.NotError(t, err)
		props[i] = prop
	}

	var resp models.Contact
	TestGet(t, "contacts/"+elem.ID, &resp)

	assert.Equals(t, elem.ID, resp.ID)
	assert.Equals(t, elem.FirstName, resp.FirstName)
	assert.Equals(t, elem.LastName, resp.LastName)
	assert.Equals(t, 3, len(resp.Props))

	for i, v := range props {
		assert.Equals(t, v.ID, resp.Props[i].ID)
		assert.Equals(t, v.Type, resp.Props[i].Type)
		assert.Equals(t, v.Name, resp.Props[i].Name)
		assert.Equals(t, v.Value, resp.Props[i].Value)
		assert.Equals(t, v.Order, resp.Props[i].Order)
	}
}

func TestSearchContact(t *testing.T) {
	ReLogin(t)
	elem, err := factories.CreateContact(models.Contact{
		UserID:   AuthUser.ID,
		DOBMonth: 1,
		DOBDay:   1,
	})
	assert.NotError(t, err)

	props := make([]models.ContactProperty, 3)
	for i := range props {
		prop, err := factories.CreateContactProperty(models.ContactProperty{
			ContactID: elem.ID,
			Name:      fmt.Sprintf("Prop %d", i+1),
			Order:     i + 1,
		})
		assert.NotError(t, err)
		props[i] = prop
	}

	_, err = factories.CreateContact(models.Contact{
		UserID:   AuthUser.ID,
		DOBMonth: 1,
		DOBDay:   1,
	})
	assert.NotError(t, err)

	deletedAt := time.Now()
	elem3, err := factories.CreateContact(models.Contact{
		UserID:    AuthUser.ID,
		DeletedAt: &deletedAt,
		DOBMonth:  1,
		DOBDay:    1,
	})
	assert.NotError(t, err)

	// should get 2 contacts
	var req requests.SearchContact
	var resp responses.SearchContact
	TestSearch(t, "contacts", req, &resp)
	assert.Equals(t, 2, len(resp.Data))

	// should get 1 deleted contact
	req = requests.SearchContact{Trashed: true}
	TestSearch(t, "contacts", req, &resp)
	assert.Equals(t, 1, len(resp.Data))
	assert.Equals(t, elem3.ID, resp.Data[0].ID)

	// should get 1 filtered contact
	req = requests.SearchContact{Query: props[0].Value}
	TestSearch(t, "contacts", req, &resp)
	assert.Equals(t, 1, len(resp.Data))
	assert.Equals(t, elem.ID, resp.Data[0].ID)
}

func TestAddContactPhoto(t *testing.T) {
	Login(t)

	picData, err := PictureBase64()
	assert.NotError(t, err)

	req := requests.FileBase64{
		Name:        "My pic",
		Description: "It is not my pic, sorry",
		Filename:    utils.RandomString() + ".png",
		Data:        picData,
	}

	elem, err := factories.CreateContact(models.Contact{
		UserID: AuthUser.ID,
	})
	assert.NotError(t, err)

	var resp models.File
	TestUpdate(t, "contacts/"+elem.ID+"/photo", req, &resp)

	assert.DatabaseHas(t, "files", utils.M{
		"id":          resp.ID,
		"user_id":     AuthUser.ID,
		"name":        req.Name,
		"description": req.Description,
		"type":        fs.FileTypeImage,
	})
	assert.DatabaseHas(t, "contacts", utils.M{
		"id":       elem.ID,
		"photo_id": resp.ID,
	})

	fileEl, err := file.GetByID(resp.ID)
	assert.NotError(t, err)

	assert.FileExists(t, fileEl.Path)

	// cleanup
	err = fs.Delete(fileEl.Path)
	assert.NotError(t, err)

	assert.FileNotExists(t, fileEl.Path)
}

func TestDeleteContactPhoto(t *testing.T) {
	Login(t)

	fileEl, err := factories.CreateFile(models.File{
		UserID: AuthUser.ID,
	})
	assert.NotError(t, err)

	elem, err := factories.CreateContact(models.Contact{
		UserID:  AuthUser.ID,
		PhotoID: fileEl.ID,
	})
	assert.NotError(t, err)

	var resp responses.Success
	TestDelete(t, "contacts/"+elem.ID+"/photo", &resp)

	assert.DatabaseMissing(t, "files", utils.M{
		"id": fileEl.ID,
	})
	assert.DatabaseHas(t, "contacts", utils.M{
		"id":       elem.ID,
		"photo_id": nil,
	})

	assert.FileNotExists(t, fileEl.Path)
}

func TestGetContactPhoto(t *testing.T) {
	Login(t)

	fileEl, err := factories.CreateFile(models.File{
		UserID: AuthUser.ID,
	})
	assert.NotError(t, err)

	elem, err := factories.CreateContact(models.Contact{
		UserID:  AuthUser.ID,
		PhotoID: fileEl.ID,
	})
	assert.NotError(t, err)

	resp := TestGet(t, "contacts/"+elem.ID+"/photo", nil)

	assert.Equals(t, fileEl.Bytes, resp.Body.Bytes())

	// cleanup
	err = fs.Delete(fileEl.Path)
	assert.NotError(t, err)
	assert.FileNotExists(t, fileEl.Path)
}

func TestGetContactPhotoPreview(t *testing.T) {
	Login(t)

	fileEl, err := factories.CreateFile(models.File{
		UserID: AuthUser.ID,
	})
	assert.NotError(t, err)

	elem, err := factories.CreateContact(models.Contact{
		UserID:  AuthUser.ID,
		PhotoID: fileEl.ID,
	})
	assert.NotError(t, err)

	TestGet(t, "contacts/"+elem.ID+"/photo/preview", nil)

	fileEl, err = file.GetByID(fileEl.ID)
	assert.NotError(t, err)

	assert.FileExists(t, fileEl.PreviewPath)

	// cleanup
	err = fs.Delete(fileEl.Path)
	assert.NotError(t, err)
	assert.FileNotExists(t, fileEl.Path)

	err = fs.Delete(fileEl.PreviewPath)
	assert.NotError(t, err)
	assert.FileNotExists(t, fileEl.PreviewPath)
}
