package workers

import (
	birthdaynotifier "github.com/vovainside/vobook/workers/birthday_notifier"
)

func Start(exit <-chan bool) {
	birthdaynotifier.Start(exit)
	//testworker.Start(exit)
}
