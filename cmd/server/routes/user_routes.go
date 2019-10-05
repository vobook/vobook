package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vovainside/vobook/cmd/server/handlers"
)

func userRoutes(r *gin.RouterGroup) {
	r.GET("user/", handlers.GetAuthUser)
}
