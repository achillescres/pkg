package mylogging

import (
	"github.com/sirupsen/logrus"
	"os"
)

func ConfigureLogrus() {
	logrus.SetReportCaller(true)
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.TraceLevel)
}

func GetLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetReportCaller(true)
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.TraceLevel)

	return logger
}
