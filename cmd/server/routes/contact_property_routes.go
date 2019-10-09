package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/vovainside/vobook/cmd/server/handlers"
)

func contactPropertyRoutes(r *gin.RouterGroup) {
	contacts := r.Group("contact-properties")
	{
		contacts.PUT("/:id/", handlers.UpdateContactProperty)
	}
}
