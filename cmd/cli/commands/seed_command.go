package commands

import (
	"sync"

	contactpropertytype "github.com/vovainside/vobook/enum/contact_property_type"

	fake "github.com/brianvoe/gofakeit"
	"github.com/vovainside/vobook/database/factories"
	"github.com/vovainside/vobook/database/models"
	"github.com/vovainside/vobook/utils"
)

func init() {
	add("seed", "sd", command{
		handler: seed,
		help:    `Seeds database with test data. Should be run on empty database`,
	})
}

func seed(args ...string) (err error) {
	// drop all data in db
	truncate := []string{
		"auth_tokens",
		"contact_properties",
		"birthday_notification_logs",
		"email_verifications",
		"password_resets",
		"contacts",
		"users",
	}
	for _, table := range truncate {
		_, err = DB.Exec(`TRUNCATE ` + table + " CASCADE;")
		if err != nil {
			return
		}
	}

	// create users
	passwordHash, err := utils.HashPassword("test")
	if err != nil {
		return
	}

	// Users
	users := []models.User{
		{
			FirstName:     "John",
			LastName:      "Snow",
			Email:         "test@test.me",
			EmailVerified: true,
		},
		{
			FirstName:     "Vladimir",
			LastName:      "Ognev",
			Email:         "vo@proj.run",
			EmailVerified: true,
		},
	}

	for i := range users {
		users[i].Password = passwordHash
	}

	err = DB.Insert(&users)
	if err != nil {
		return
	}
	println(len(users), "users created")

	// create contacts
	var wg sync.WaitGroup

	for _, u := range users {
		for i := 0; i < 30; i++ {
			go func(user models.User) {
				wg.Add(1)
				defer wg.Done()

				elem := models.Contact{
					UserID:    user.ID,
					DeletedAt: nil,
				}
				contact, err := factories.CreateContact(elem)
				if err != nil {
					println(err.Error())
					return
				}

				props := make([]models.ContactProperty, 0)
				for i, v := range contactpropertytype.All {
					if fake.Bool() {
						continue
					}

					prop := models.ContactProperty{
						ContactID: contact.ID,
						Order:     i,
						Name:      v.String(),
						Type:      v,
					}

					switch v {
					case contactpropertytype.PersonalPhone, contactpropertytype.WorkPhone, contactpropertytype.Phone:
						prop.Value = fake.Phone()
					case contactpropertytype.WorkEmail, contactpropertytype.PersonalEmail, contactpropertytype.Email:
						prop.Value = fake.Email()
					}

					prop, err = factories.MakeContactProperty(prop)
					if err != nil {
						println(err.Error())
					}

					props = append(props, prop)
				}
				err = DB.Insert(&props)
				if err != nil {
					println(err.Error())
				}
			}(u)
		}
	}
	wg.Wait()

	return
}
