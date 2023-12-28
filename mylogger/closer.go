package mylogging

import (
	"github.com/sirupsen/logrus"
)

type Closer func(*logrus.Entry) error
