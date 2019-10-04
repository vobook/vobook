package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vovainside/vobook/cmd/server/handlers"
)

func emailVerificationRoutes(r *gin.Engine) {
	r.GET("verify-email/:id", handlers.VerifyEmail)
}
