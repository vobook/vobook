package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/vovainside/vobook/cmd/server/handlers"
	"github.com/vovainside/vobook/cmd/server/middlewares"
)

func contactPropertyRoutes(r *gin.RouterGroup) {
	all := r.Group("contact-properties")
	{
		one := all.Group("/:id/")
		one.Use(middlewares.ContactPropertyMiddleware)
		{
			one.PUT("/", handlers.UpdateContactProperty)
		}
	}

	r.PUT("/trash-contact-properties/", handlers.TrashContactProperties)
	r.PUT("/restore-contact-properties/", handlers.RestoreContactProperties)
	r.PUT("/delete-contact-properties/", handlers.DeleteContactProperties)
	r.PUT("/reorder-contact-properties/", handlers.ReorderContactProperties)
}
