package handlers

import (
	"net/http"
	"vobook/cmd/server/errors"
	contactproperty "vobook/domain/contact_property"
	"vobook/domain/file"
	"vobook/services/fs"

	"github.com/gin-gonic/gin"

	"vobook/cmd/server/requests"
	"vobook/cmd/server/responses"
	"vobook/database/models"
	"vobook/domain/contact"
)

func CreateContact(c *gin.Context) {
	var req requests.CreateContact
	if !bindJSON(c, &req) {
		return
	}

	elem, err := req.ToModel()
	if err != nil {
		Abort(c, err)
		return
	}
	elem.UserID = AuthUser(c).ID
	err = contact.Create(elem)
	if err != nil {
		Abort(c, err)
		return
	}

	c.JSON(http.StatusCreated, elem)
}

func SearchContacts(c *gin.Context) {
	var req requests.SearchContact
	if !bindQuery(c, &req) {
		return
	}

	data, count, err := contact.Search(AuthUser(c).ID, req)
	if err != nil {
		Abort(c, err)
		return
	}

	resp := responses.SearchContact{
		Data:  data,
		Count: count,
	}

	c.JSON(http.StatusOK, resp)
}

func UpdateContact(c *gin.Context) {
	var req requests.UpdateContact
	if !bindJSON(c, &req) {
		return
	}

	elem := getContactFromRequest(c)
	err := req.ToModel(&elem)
	if err != nil {
		Abort(c, err)
		return
	}
	err = contact.Update(&elem)
	if err != nil {
		Abort(c, err)
		return
	}

	if len(elem.Props) > 0 {
		err = contactproperty.DeleteByContact(elem.ID)
		if err != nil {
			Abort(c, err)
			return
		}
		for i := range elem.Props {
			elem.Props[i].ContactID = elem.ID
		}
		err = contactproperty.CreateMany(&elem.Props)
		if err != nil {
			Abort(c, err)
			return
		}
	}

	c.JSON(http.StatusOK, responses.OK("Saved"))
}

func TrashContacts(c *gin.Context) {
	var ids requests.IDs
	if !bindJSON(c, &ids) {
		return
	}

	err := contact.Trash(AuthUser(c).ID, ids...)
	if err != nil {
		Abort(c, err)
		return
	}

	c.JSON(http.StatusOK, responses.OK("Success"))
}

func RestoreContacts(c *gin.Context) {
	var ids requests.IDs
	if !bindJSON(c, &ids) {
		return
	}

	err := contact.Restore(AuthUser(c).ID, ids...)
	if err != nil {
		Abort(c, err)
		return
	}

	c.JSON(http.StatusOK, responses.OK("Success"))
}

func DeleteContacts(c *gin.Context) {
	var ids requests.IDs
	if !bindJSON(c, &ids) {
		return
	}

	err := contact.Delete(AuthUser(c).ID, ids...)
	if err != nil {
		Abort(c, err)
		return
	}

	c.JSON(http.StatusOK, responses.OK("Success"))
}

func GetContact(c *gin.Context) {
	elem := getContactFromRequest(c)
	props, err := contact.Props(elem.ID)
	if err != nil {
		Abort(c, err)
		return
	}

	elem.Props = props
	c.JSON(http.StatusOK, elem)
}

func AddContactPhoto(c *gin.Context) {
	contactEl := getContactFromRequest(c)

	var req requests.FileBase64
	if !bindJSON(c, &req) {
		return
	}

	photoEl, err := req.ToModel()
	if err != nil {
		Abort(c, err)
		return
	}
	photoEl.UserID = AuthUser(c).ID

	err = contact.AddPhoto(contactEl.ID, photoEl)
	if err != nil {
		Abort(c, err)
		return
	}

	c.JSON(http.StatusOK, photoEl)
}

func GetContactPhoto(c *gin.Context) {
	elem := getContactFromRequest(c)

	if elem.PhotoID == "" {
		Abort(c, errors.ContactPhotoNotExists)
		return
	}

	fileEl, err := file.GetByID(elem.PhotoID)
	if err != nil {
		Abort(c, err)
		return
	}

	c.File(fs.FullPath(fileEl.Path))
}

func GetContactPhotoPreview(c *gin.Context) {
	elem := getContactFromRequest(c)

	if elem.PhotoID == "" {
		Abort(c, errors.ContactPhotoNotExists)
		return
	}

	fileEl, err := file.GetByID(elem.PhotoID)
	if err != nil {
		Abort(c, err)
		return
	}

	if fileEl.PreviewPath == "" {
		err = file.CreatePreview(&fileEl)
		if err != nil {
			Abort(c, err)
			return
		}
	}

	c.File(fs.FullPath(fileEl.PreviewPath))
}

func DeleteContactPhoto(c *gin.Context) {
	elem := getContactFromRequest(c)

	err := contact.DeletePhoto(elem)
	if err != nil {
		Abort(c, err)
		return
	}

	c.JSON(http.StatusOK, responses.OK("Success"))
}

func getContactFromRequest(c *gin.Context) models.Contact {
	elem := c.MustGet("contact")
	return elem.(models.Contact)
}
