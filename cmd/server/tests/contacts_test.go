package tests

import (
	"testing"
	"time"

	fake "github.com/brianvoe/gofakeit"
	"github.com/vovainside/vobook/cmd/server/requests"
)

func TestCreateContact(t *testing.T) {
	req := requests.CreateContact{
		Name:       fake.Name(),
		FirstName:  fake.FirstName(),
		LastName:   fake.LastName(),
		Birthday:   fake.DateRange(time.Now().AddDate(-100, 0, 0), time.Now()),
		Properties: []requests.CreateContactProperty{
			//Type:  contactproperty.TypeEmail,
			//Value: fake.Email(),
		},
	}

	//var resp models.User
	//POST(t, Request{
	//	Path:         "register-user",
	//	Body:         req,
	//	AssertStatus: http.StatusOK,
	//	BindResponse: &resp,
	//	IsPublic:     true,
	//})

}
