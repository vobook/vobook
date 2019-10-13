package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/vovainside/vobook/cmd/server/requests"
	"github.com/vovainside/vobook/cmd/server/responses"
	"github.com/vovainside/vobook/database/models"
	"github.com/vovainside/vobook/domain/contact"
)

func CreateContact(c *gin.Context) {
	var req requests.CreateContact
	if !bindJSON(c, &req) {
		return
	}

	elem := req.ToModel()
	elem.UserID = authUser(c).ID
	err := contact.Create(elem)
	if err != nil {
		Abort(c, err)
		return
	}

	c.JSON(http.StatusCreated, elem)
}

func SearchContact(c *gin.Context) {
	var req requests.SearchContact
	if !bindQuery(c, &req) {
		return
	}

	data, count, err := contact.Search(authUser(c).ID, elem)
	if err != nil {
		Abort(c, err)
		return
	}

	resp := responses.SearchContact{
		Data:  data,
		Count: count,
	}

	c.JSON(http.StatusCreated, resp)
}

func UpdateContact(c *gin.Context) {
	var req requests.UpdateContact
	if !bindJSON(c, &req) {
		return
	}

	elem := getContactFromRequest(c)
	req.ToModel(&elem)
	err := contact.Update(&elem)
	if err != nil {
		Abort(c, err)
		return
	}

	c.JSON(http.StatusOK, responses.OK("Saved"))
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

func getContactFromRequest(c *gin.Context) models.Contact {
	elem := c.MustGet("contact")
	return elem.(models.Contact)
}
