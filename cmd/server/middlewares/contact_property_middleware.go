package middlewares

import (
	"github.com/gin-gonic/gin"

	"github.com/vovainside/vobook/cmd/server/errors"
	"github.com/vovainside/vobook/cmd/server/handlers"
	"github.com/vovainside/vobook/domain/contact_property"
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
