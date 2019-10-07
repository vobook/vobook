package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vovainside/vobook/cmd/server/middlewares"
	"github.com/vovainside/vobook/config"
)

func Register(r *gin.Engine) {
	conf := config.Get()

	api := r.Group(conf.ApiBasePath)
	api.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, conf.App.Name+" "+conf.App.Version+"."+conf.App.Build+"."+conf.App.Env)
	})

	api.Use(
		middlewares.ClientID,
	)

	publicRoutes(api)

	api.Use(
		middlewares.TokenAuth,
	)

	apply(api,
		userRoutes,
		contactRoutes,
	)
}

func apply(rg *gin.RouterGroup, routesFn ...func(*gin.RouterGroup)) {
	for _, fn := range routesFn {
		fn(rg)
	}
}
