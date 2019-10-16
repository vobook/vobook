package main

import (
	"os"
	"os/signal"

	log "github.com/sirupsen/logrus"

	"github.com/vovainside/vobook/config"
	"github.com/vovainside/vobook/logger"
	"github.com/vovainside/vobook/services/mail"
	"github.com/vovainside/vobook/workers"
)

func main() {
	conf := config.Get()
	logger.Setup()

	mail.InitDrivers()

	quit := make(chan os.Signal)
	exit := make(chan bool)
	signal.Notify(quit, os.Interrupt)

	workers.Start(exit)

	go func() {
		<-quit
		close(exit)
	}()

	<-exit
	log.Println(conf.App.Name + "." + conf.App.Env + " workers stopped")
}
