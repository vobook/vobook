package middlewares

import (
	"github.com/gin-gonic/gin"

	"github.com/vovainside/vobook/cmd/server/handlers"
	"github.com/vovainside/vobook/domain/contact"
)

func ContactMiddleware(c *gin.Context) {
	elem, err := contact.Find(c.Param("id"))
	if err != nil {
		handlers.Abort(c, err)
		return
	}

	c.Set("contact", elem)
	c.Next()
}
