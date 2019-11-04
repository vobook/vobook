package contactpropertytype

import (
	"fmt"
	"math/rand"

	"github.com/vovainside/vobook/cmd/server/errors"
)

type Type int

const (
	_             Type = iota
	Other              // 1
	PersonalPhone      // 2
	WorkPhone          // 3
	PersonalEmail      // 4
	WorkEmail          // 5
	Phone              // 6
	Email              // 7
	Address            // 8
	Facebook           // 9
	Twitter            // 10
	Github             // 11
	WhatsApp           // 12
	Telegram           // 13
	Link               // 15
	Note               // 16
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

type TypeModel struct {
	Type Type   `json:"type"`
	Name string `json:"name"`
}

func Random() Type {
	return All[rand.Intn(len(All)-1)+1]
}
