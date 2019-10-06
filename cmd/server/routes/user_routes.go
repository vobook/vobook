package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vovainside/vobook/cmd/server/handlers"
)

func userRoutes(r *gin.RouterGroup) {
	r.GET("user/", handlers.GetAuthUser)
	r.PUT("user/", handlers.UpdateAuthUser)
	r.POST("change-password/", handlers.ChangePassword)
	r.POST("change-email/", handlers.ChangeEmail)
}
