package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/vovainside/vobook/cmd/server/requests"
	"github.com/vovainside/vobook/cmd/server/responses"
	"github.com/vovainside/vobook/domain/contact_property"
)

func UpdateContactProperty(c *gin.Context) {
	var req requests.UpdateContactProperty
	if !bindJSON(c, &req) {
		return
	}

	id := c.Param("id")
	elem, err := contactproperty.Find(id)
	if err != nil {
		Abort(c, err)
		return
	}

	err = contactproperty.Update(&elem)
	if err != nil {
		Abort(c, err)
		return
	}

	c.JSON(http.StatusOK, responses.OK("Saved"))
}
