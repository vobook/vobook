package logger

import (
	"io"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/vovainside/vobook/config"
	"gopkg.in/natefinch/lumberjack.v2"
)

func Setup() {
	log.SetFormatter(&Formatter{})
	log.SetLevel(log.DebugLevel)

	filename := config.Get().LogsFilePath
	if filename == "" {
		filename = "vobook.log"
	}

	fileWriter := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    5, // megabytes
		MaxBackups: 3,
		MaxAge:     5,    //days
		Compress:   true, // disabled by default
	}

	log.SetOutput(io.MultiWriter(os.Stdout, fileWriter))
}
