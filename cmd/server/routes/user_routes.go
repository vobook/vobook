package routes

import (
	"vobook/cmd/server/handlers"

	"github.com/gin-gonic/gin"
)

func userRoutes(r *gin.RouterGroup) {
	r.GET("user/", handlers.GetAuthUser)
	r.PUT("user/", handlers.UpdateAuthUser)
	r.PUT("user/delete/", handlers.DeleteAuthUser)
	//r.PUT("user/restore/", handlers.RestoreAuthUser) // TODO
	//r.PUT("user/erase/", handlers.EraseAuthUser) // TODO
	r.POST("change-password/", handlers.ChangePassword)
	r.POST("change-email/", handlers.ChangeEmail)
}
