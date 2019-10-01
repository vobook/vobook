package routes

import (
	"net/http"

	"github.com/vovainside/vobook/cmd/server/handlers"
)

func init() {
	Add(
		Route{
			Method:  http.MethodGet,
			Path:    "users",
			Handler: handlers.SearchUsers,
		},
		Route{
			Method:  http.MethodPost,
			Path:    "users/register",
			Handler: handlers.RegisterUser,
		},
		Route{
			Method:  http.MethodGet,
			Path:    "users/:id",
			Handler: handlers.GetUserByID,
		},
	)
}
