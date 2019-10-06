package mail

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"sync"

	"github.com/davecgh/go-spew/spew"
	"github.com/vovainside/vobook/config"
)

func NewTestDriver() Driver {
	return TestDriver{}
}

type TestDriver struct{}

func (drv TestDriver) Send(msg Message) (err error) {
	mBytes, err := json.MarshalIndent(msg, "", "    ")
	if err != nil {
		return
	}

	stub := config.Get().Mail.Stub
	err = os.MkdirAll(stub, 0755)
	if err != nil {
		spew.Dump(err)
		return
	}

	err = ioutil.WriteFile(path.Join(stub, msg.To[0]+".json"), mBytes, 0755)
	if err != nil {
		return
	}

	TestRepo.Store(msg.To[0], msg)
	return
}

type testMails struct {
	sync.Map
}

func (m testMails) GetMail(to string) (msg *Message) {
	v, ok := m.Load(to)
	if !ok {
		return nil
	}

	message := v.(Message)
	return &message
}

var TestRepo = testMails{}
