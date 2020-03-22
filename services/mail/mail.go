package mail

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"path"

	"vobook/config"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

const TestDriverName = "test"

var Drivers map[string]Driver

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

type Replace map[string]string

func InitDrivers() {
	Drivers = map[string]Driver{
		"go-mail":      NewGoMailDriver(),
		TestDriverName: NewTestDriver(),
	}
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

func LoadTemplate(name string, replace Replace) (m Message, err error) {
	fname := path.Join(config.Get().Mail.Templates, name+".yaml")
	data, err := ioutil.ReadFile(fname)
	if err != nil {
		return
	}

	data = replacePlaceholders(data, replace)
	err = yaml.Unmarshal(data, &m)
	return
}

func SendTemplate(to, template string, replace Replace) (err error) {
	m, err := LoadTemplate(template, replace)
	if err != nil {
		return
	}

	m.To = []string{to}
	return Send(m)
}

func replacePlaceholders(data []byte, replace Replace) []byte {
	globals := Replace{
		"app-name": config.Get().App.Name,
	}

	if replace == nil {
		replace = globals
		goto replace
	}

	for k, v := range globals {
		_, ok := replace[k]
		if !ok {
			replace[k] = v
		}
	}

replace:
	for k, v := range replace {
		data = bytes.ReplaceAll(data, []byte("{"+k+"}"), []byte(v))
	}

	return data
}
