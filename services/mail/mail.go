package mail

import (
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/vovainside/vobook/config"
)

const TestDriverName = "test"

type Message struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Cc      []string `json:"cc"`
	Subject string   `json:"subject"`
	Body    string   `json:"body"`
}

type Driver interface {
	Send(m Message) error
}

var Drivers = map[string]Driver{
	"go-mail":      NewGoMailDriver(),
	TestDriverName: NewTestDriver(),
}

var ErrToAddrNotSet = errors.New("to address not set")

func Send(m Message) (err error) {
	name := config.Get().Mail.Driver
	drv, ok := Drivers[name]
	if !ok {
		panic(fmt.Sprintf("email: unknown driver '%s'", name))
	}

	if m.From == "" {
		m.From = config.Get().Mail.From
	}

	if len(m.To) == 0 {
		err = ErrToAddrNotSet
		return
	}

	if name != TestDriverName {
		go func() {
			err = drv.Send(m)
			if err != nil {
				log.Error(err)
			}
		}()
		return
	}

	return drv.Send(m)
}

func SendSimple(to, subject, body string) (err error) {
	return Send(Message{
		To:      []string{to},
		Subject: subject,
		Body:    body,
	})
}
