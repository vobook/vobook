package contactpropertytype

import (
	"fmt"
	"math/rand"

	"github.com/vovainside/vobook/cmd/server/errors"
)

type Type int

const (
	Other         Type = iota
	PersonalPhone      // 1
	WorkPhone          // 2
	PersonalEmail      // 3
	WorkEmail          // 4
	Phone              // 5
	Email              // 6
	Address            // 7
	Facebook           // 8
	Twitter            // 9
	Github             // 10
	WhatsApp           // 11
	Telegram           // 12
	Link               // 13
	Note               // 14
)

func (t Type) String() string {
	switch t {
	case Other:
		return "Other"
	case PersonalPhone:
		return "Personal phone"
	case WorkPhone:
		return "Work phone"
	case Phone:
		return "Phone"
	case Email:
		return "Email"
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
	case Link:
		return "Link"
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
