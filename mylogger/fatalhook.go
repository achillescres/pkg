package mylogging

import "github.com/sirupsen/logrus"

type CloseOnFatal struct {
	closer Closer
}

func NewCloseOnFatal(closer Closer) *CloseOnFatal {
	return &CloseOnFatal{closer: closer}
}

func (CloseOnFatal) Levels() []logrus.Level {
	return []logrus.Level{logrus.FatalLevel}
}

func (hook CloseOnFatal) Fire(e *logrus.Entry) error {
	e.Infoln("Hook CloseOnFatal!")
	err := hook.closer(e)
	e.Infoln("Hook CloseOnFatal success!")
	return err
}
