package middlewares

import (
	"github.com/gin-gonic/gin"

	"vobook/cmd/server/errors"
	"vobook/cmd/server/handlers"
	"vobook/domain/contact"
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
