package tests

import (
	"testing"

	fake "github.com/brianvoe/gofakeit"

	"github.com/vovainside/vobook/cmd/server/requests"
	"github.com/vovainside/vobook/cmd/server/responses"
	"github.com/vovainside/vobook/database/factories"
	. "github.com/vovainside/vobook/tests/apitest"
	"github.com/vovainside/vobook/tests/assert"
	"github.com/vovainside/vobook/utils"
)

func TestUpdateContactProperty(t *testing.T) {
	prop, err := factories.CreateContactProperty()
	if err != nil {
		return
	}
	name := fake.Name()
	value := fake.Name()
	req := requests.UpdateContactProperty{
		Name:  &name,
		Value: &value,
	}

	var resp responses.Success
	TestUpdate(t, "contact-properties/"+prop.ID, req, &resp)

	assert.DatabaseHas(t, "contact_properties", utils.M{
		"id":    prop.ID,
		"type":  prop.Type,
		"name":  name,
		"value": value,
	})
}
