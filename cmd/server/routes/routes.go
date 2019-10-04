package routes

import (
	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine) {
	userRoutes(r)
	emailVerificationRoutes(r)
}
