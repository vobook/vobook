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

		one := all.Group(":id")
		one.Use(middlewares.ContactMiddleware)
		{
			one.PUT("/", handlers.UpdateContact)
		}

	}
}
