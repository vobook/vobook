package responses

import "github.com/vovainside/vobook/database/models"

type SearchContact struct {
	Data  []models.Contact `json:"data"`
	Count int              `json:"count"`
}
