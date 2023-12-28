package mylogging

import "github.com/sirupsen/logrus"

type CloseOnFatal struct {
	exec func(*logrus.Entry) error
}

func (*CloseOnFatal) Levels() []logrus.Level {
	return []logrus.Level{logrus.FatalLevel}
}

func (hook *CloseOnFatal) Fire(e *logrus.Entry) error {
	e.Infoln("Hook CloseOnFatal!")
	err := hook.exec(e)
	e.Infoln("Hook CloseOnFatal success!")
	return err
}
