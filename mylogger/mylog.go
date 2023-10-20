package mylogging

import (
	customLog "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
	"os"
)

func ConfigureLogrus() {
	logrus.SetFormatter(&customLog.Formatter{
		FieldsOrder: []string{"component"},
	})
	logrus.SetReportCaller(true)
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.TraceLevel)
}

func GetLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&customLog.Formatter{
		FieldsOrder: []string{"component", "category"},
	})
	logger.SetReportCaller(true)
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.TraceLevel)

	return logger
}
