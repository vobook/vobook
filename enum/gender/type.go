package gender

import (
	"fmt"
	"math/rand"

	"github.com/vovainside/vobook/cmd/server/errors"
)

type Type int

const (
	Unknown Type = iota
	Male
	Female
	Other
)

func (t Type) Name() string {
	switch t {
	case Unknown:
		return "Unknown"
	case Male:
		return "Male"
	case Female:
		return "Female"
	case Other:
		return "Other"
	}

	panic(fmt.Sprintf("unknown gender %d", t))
}

func (t Type) Validate() (err error) {
	for _, v := range Types {
		if v == t {
			return nil
		}
	}

	return errors.InvalidGender
}

var Types = []Type{
	Unknown,
	Male,
	Female,
	Other,
}

func RandomType() Type {
	return Types[rand.Intn(len(Types)-1)]
}
