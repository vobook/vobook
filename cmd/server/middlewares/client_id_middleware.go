package middlewares

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vovainside/vobook/cmd/server/errors"
	"github.com/vovainside/vobook/cmd/server/handlers"
	"github.com/vovainside/vobook/database/models"
)

func ClientID(c *gin.Context) {
	client := c.Request.Header.Get("X-Client")
	if client == "" {
		handlers.Abort(c, errors.InvalidAppClient)
		return
	}

	clientID, err := strconv.Atoi(client)
	if err != nil {
		handlers.Abort(c, errors.InvalidAppClient)
		return
	}

	for _, elem := range models.Clients {
		if elem == models.ClientID(clientID) {
			c.Set("clientID", clientID)
			c.Next()
			return
		}
	}

	handlers.Abort(c, errors.InvalidAppClient)
	return
}
