// This is my proof of concept for base worker
// Worker should process some some task evert X interval
package testworker

import (
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	checkInterval = 1 * time.Minute
)

func Start(exit <-chan bool) {
	go worker(exit)
}

func worker(exit <-chan bool) {
	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()

	log.Println("Test worker started")
loop:
	for {
		select {
		case <-exit:
			log.Println("Test worker stopped")
			break loop
		case <-ticker.C:
			go work()
		}
	}
}

func work() {
	log.Println("[test worker] hey, im working!")
}
