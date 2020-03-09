package routes

import (
	"vobook/cmd/server/handlers"

	"github.com/gin-gonic/gin"
)

func publicRoutes(r *gin.RouterGroup) {
	r.POST("register-user/", handlers.RegisterUser)
	r.POST("login/", handlers.Login)
	r.GET("verify-email/:token/", handlers.VerifyEmail)
	r.POST("reset-password/", handlers.ResetPasswordStart)
	r.POST("reset-password/:token/", handlers.ResetPasswordCheckToken)
	r.PUT("reset-password/", handlers.ResetPassword)
}
