package middlewares

import (
	"github.com/gin-gonic/gin"

	"vobook/cmd/server/errors"
	"vobook/cmd/server/handlers"
	contactproperty "vobook/domain/contact_property"
)

func ContactPropertyMiddleware(c *gin.Context) {
	prop, err := contactproperty.Find(c.Param("id"))
	if err != nil {
		handlers.Abort(c, err)
		return
	}

	if prop.Contact.UserID != handlers.AuthUser(c).ID {
		handlers.Abort(c, errors.ContactNotFound)
		return
	}

	c.Set("contact-property", prop)
	c.Next()
}
