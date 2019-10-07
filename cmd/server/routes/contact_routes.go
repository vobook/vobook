package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vovainside/vobook/cmd/server/handlers"
)

func contactRoutes(r *gin.RouterGroup) {
	contacts := r.Group("contacts")
	{
		contacts.POST("/", handlers.CreateContact)
	}
}
