package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/vovainside/vobook/cmd/server/handlers"
	"github.com/vovainside/vobook/cmd/server/middlewares"
)

func contactRoutes(r *gin.RouterGroup) {
	all := r.Group("contacts")
	{
		all.POST("/", handlers.CreateContact)
		all.GET("/", handlers.SearchContacts)

		one := all.Group(":id")
		one.Use(middlewares.ContactMiddleware)
		{
			one.PUT("/", handlers.UpdateContact)
			one.GET("/", handlers.GetContact)
		}

	}

	// because gin can't have routes that will match by willcard
	// (panic: wildcard route ':id' conflicts with existing children in path)
	// i'll gonna use verbs to get around this limitation
	r.PUT("/trash-contacts/", handlers.TrashContacts)
}
