package workers

import (
	"github.com/vovainside/vobook/workers/birthday_notifier"
	"github.com/vovainside/vobook/workers/test_worker"
)

func Start(exit <-chan bool) {
	birthdaynotifier.Start(exit)
	testworker.Start(exit)
}
