package routes

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/vovainside/vobook/cmd/server/middlewares"
	"github.com/vovainside/vobook/config"
)

func Register(r *gin.Engine) {
	conf := config.Get()

	r.Use(corsConfig())

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
		contactPropertyRoutes,
	)
}

func apply(rg *gin.RouterGroup, routesFn ...func(*gin.RouterGroup)) {
	for _, fn := range routesFn {
		fn(rg)
	}
}

func corsConfig() gin.HandlerFunc {
	conf := cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // TODO move to conf
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "X-Client"},
		AllowCredentials: false,
		MaxAge:           24 * time.Hour,
	}

	return cors.New(conf)
}
