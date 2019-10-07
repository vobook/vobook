package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vovainside/vobook/cmd/server/requests"
	"github.com/vovainside/vobook/cmd/server/responses"
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

	c.JSON(http.StatusCreated, responses.OK("New contact created"))
}
