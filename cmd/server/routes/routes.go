package routes

import (
	"testing"

	"github.com/gin-gonic/gin"
)

var routes []Route

type Route struct {
	Name        string
	Description string
	Method      string
	Path        string
	Handler     gin.HandlerFunc
	MiddleWares []gin.HandlerFunc
	Test        func(t *testing.T)
}

func Add(r ...Route) {
	routes = append(routes, r...)
}

func Register(router *gin.Engine) {
	for _, r := range routes {
		if len(r.MiddleWares) > 0 {
			group := router.Group("/", r.MiddleWares...)
			group.Handle(r.Method, r.Path, r.Handler)
			continue
		}
		router.Handle(r.Method, r.Path, r.Handler)

	}

	routes = make([]Route, 0)
}
