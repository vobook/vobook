package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/vovainside/vobook/cmd/server/requests"
	"github.com/vovainside/vobook/cmd/server/responses"
	"github.com/vovainside/vobook/database/models"
	"github.com/vovainside/vobook/domain/contact_property"
)

func UpdateContactProperty(c *gin.Context) {
	var req requests.UpdateContactProperty
	if !bindJSON(c, &req) {
		return
	}

	elem := getContactPropertyFromRequest(c)
	req.ToModel(&elem)

	err := contactproperty.Update(&elem)
	if err != nil {
		Abort(c, err)
		return
	}

	c.JSON(http.StatusOK, responses.OK("Saved"))
}

func TrashContactProperties(c *gin.Context) {
	var req requests.IDs
	if !bindJSON(c, &req) {
		return
	}

	err := contactproperty.Trash(AuthUser(c).ID, req...)
	if err != nil {
		Abort(c, err)
		return
	}

	c.JSON(http.StatusOK, responses.OK("Deleted"))
}

func RestoreContactProperties(c *gin.Context) {
	var req requests.IDs
	if !bindJSON(c, &req) {
		return
	}

	err := contactproperty.Restore(AuthUser(c).ID, req...)
	if err != nil {
		Abort(c, err)
		return
	}

	c.JSON(http.StatusOK, responses.OK("Restored"))
}

func DeleteContactProperties(c *gin.Context) {
	var req requests.IDs
	if !bindJSON(c, &req) {
		return
	}

	err := contactproperty.Delete(AuthUser(c).ID, req...)
	if err != nil {
		Abort(c, err)
		return
	}

	c.JSON(http.StatusOK, responses.OK("Deleted"))
}

func getContactPropertyFromRequest(c *gin.Context) models.ContactProperty {
	elem := c.MustGet("contact-property")
	return elem.(models.ContactProperty)
}
