package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vovainside/vobook/cmd/server/handlers"
)

func publicRoutes(r *gin.RouterGroup) {
	r.POST("register-user/", handlers.RegisterUser)
	r.POST("login/", handlers.Login)
	r.GET("verify-email/:id/", handlers.VerifyEmail)
}
