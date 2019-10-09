package contactproperty

import (
	"fmt"
	"math/rand"

	"github.com/vovainside/vobook/cmd/server/errors"
)

type Type int

const (
	TypeOther Type = iota
	TypePersonalPhone
	TypeWorkPhone
	TypePersonalEmail
	TypeWorkEmail
	TypePhone
	TypeEmail
	TypeAddress
	TypeFacebook
	TypeTwitter
	TypeGithub
	TypeWhatsApp
	TypeTelegram
	TypeLink
	TypeNote
)

func (t Type) Name() string {
	switch t {
	case TypeOther:
		return "Other"
	case TypePersonalPhone:
		return "Personal phone"
	case TypeWorkPhone:
		return "Work phone"
	case TypePersonalEmail:
		return "Personal email"
	case TypeWorkEmail:
		return "Work email"
	case TypeAddress:
		return "Address"
	case TypeFacebook:
		return "Facebook"
	case TypeTwitter:
		return "Twitter"
	case TypeGithub:
		return "Github"
	case TypeWhatsApp:
		return "WhatsApp"
	case TypeTelegram:
		return "Telegram"
	case TypeNote:
		return "Note"
	}

	panic(fmt.Sprintf("unknown contact property type %d", t))
}

func (t Type) Validate() (err error) {
	for _, v := range Types {
		if v == t {
			return nil
		}
	}

	return errors.InvalidContactPropertyType
}

var Types = []Type{
	TypeOther,
	TypePersonalPhone,
	TypeWorkPhone,
	TypePersonalEmail,
	TypeWorkEmail,
	TypePhone,
	TypeEmail,
	TypeAddress,
	TypeFacebook,
	TypeTwitter,
	TypeGithub,
	TypeWhatsApp,
	TypeTelegram,
	TypeLink,
	TypeNote,
}

func RandomType() Type {
	return Types[rand.Intn(len(Types)-1)]
}
