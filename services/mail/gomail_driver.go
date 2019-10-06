package mail

import (
	"crypto/tls"

	"github.com/vovainside/vobook/config"
	"gopkg.in/gomail.v2"
)

func NewGoMailDriver() Driver {
	cm := config.Get().Mail
	d := gomail.NewDialer(cm.Host, cm.Port, cm.User, cm.Password)
	d.TLSConfig = &tls.Config{
		InsecureSkipVerify: true,
	}

	return GoMailDriver{
		Dialer: d,
	}
}

type GoMailDriver struct {
	Dialer *gomail.Dialer
}

func (drv GoMailDriver) Send(msg Message) (err error) {
	m := gomail.NewMessage()
	m.SetHeader("From", msg.From)
	m.SetHeader("To", msg.To...)
	m.SetHeader("Subject", msg.Subject)
	m.SetBody("text/html", msg.Body)

	if len(msg.Cc) > 0 {
		m.SetHeader("Cc", msg.Cc...)
	}

	err = drv.Dialer.DialAndSend(m)
	return
}
