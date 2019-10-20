package contactpropertytype

import (
	"fmt"
	"math/rand"

	"github.com/vovainside/vobook/cmd/server/errors"
)

type Type int

const (
	Other Type = iota
	PersonalPhone
	WorkPhone
	PersonalEmail
	WorkEmail
	Phone
	Email
	Address
	Facebook
	Twitter
	Github
	WhatsApp
	Telegram
	Link
	Note
)

func (t Type) String() string {
	switch t {
	case Other:
		return "Other"
	case PersonalPhone:
		return "Personal phone"
	case WorkPhone:
		return "Work phone"
	case PersonalEmail:
		return "Personal email"
	case WorkEmail:
		return "Work email"
	case Address:
		return "Address"
	case Facebook:
		return "Facebook"
	case Twitter:
		return "Twitter"
	case Github:
		return "Github"
	case WhatsApp:
		return "WhatsApp"
	case Telegram:
		return "Telegram"
	case Note:
		return "Note"
	}

	panic(fmt.Sprintf("unknown contact property type %d", t))
}

func (t Type) Validate() (err error) {
	for _, v := range All {
		if v == t {
			return nil
		}
	}

	return errors.InvalidContactPropertyType
}

var All = []Type{
	Other,
	PersonalPhone,
	WorkPhone,
	PersonalEmail,
	WorkEmail,
	Phone,
	Email,
	Address,
	Facebook,
	Twitter,
	Github,
	WhatsApp,
	Telegram,
	Link,
	Note,
}

func Random() Type {
	return All[rand.Intn(len(All)-1)]
}
