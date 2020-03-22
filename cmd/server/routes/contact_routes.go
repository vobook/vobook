package routes

import (
	"github.com/gin-gonic/gin"

	"vobook/cmd/server/handlers"
	"vobook/cmd/server/middlewares"
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
			photo := one.Group("photo")
			{
				photo.GET("/", handlers.GetContactPhoto)
				photo.GET("/preview/", handlers.GetContactPhotoPreview)
				photo.PUT("/", handlers.AddContactPhoto)
				photo.DELETE("/", handlers.DeleteContactPhoto)
			}
		}
	}

	// because gin can't have routes that will match by wildcard
	// (panic: wildcard route ':id' conflicts with existing children in path)
	// i'll gonna use verbs to get around this limitation
	r.PUT("/trash-contacts/", handlers.TrashContacts)
	r.PUT("/restore-contacts/", handlers.RestoreContacts)
	r.PUT("/delete-contacts/", handlers.DeleteContacts)
}
