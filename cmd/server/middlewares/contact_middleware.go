package middlewares

import (
	"github.com/gin-gonic/gin"

	"github.com/vovainside/vobook/cmd/server/errors"
	"github.com/vovainside/vobook/cmd/server/handlers"
	"github.com/vovainside/vobook/domain/contact"
)

func ContactMiddleware(c *gin.Context) {
	elem, err := contact.Find(c.Param("id"))
	if err != nil {
		handlers.Abort(c, err)
		return
	}

	if elem.UserID != handlers.AuthUser(c).ID {
		handlers.Abort(c, errors.ContactNotFound)
		return
	}

	c.Set("contact", elem)
	c.Next()
}
