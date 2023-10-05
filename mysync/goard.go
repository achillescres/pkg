package mysync

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"runtime"
)

var logger *logrus.Entry = logrus.NewEntry(logrus.StandardLogger())

// ChangeLogger change logger to newLogger, use nil to prevent Goard from logging
func ChangeLogger(newLogger *logrus.Entry) {
	logger = newLogger
}

func Goard(f func()) {
	go func() {
		defer func() {
			recover()
		}()
		f()
	}()
}

func GoardWithLog(f func()) {
	go func() {
		//defer func() {
		//	if r := recover(); logger != nil && r != nil {
		//		buf := make([]byte, 64<<10)
		//		buf = buf[:runtime.Stack(buf, false)]
		//		err := fmt.Errorf("errgroup: panic recovered: %s\n%s", r, buf)
		//		logger.Errorln(err)
		//	}
		//}()
		f()
	}()
}

func GoardWithLogger(f func(), logger *logrus.Logger) {
	go func() {
		//defer func() {
		//	if r := recover(); logger != nil && r != nil {
		//		buf := make([]byte, 64<<10)
		//		buf = buf[:runtime.Stack(buf, false)]
		//		err := fmt.Errorf("errgroup: panic recovered: %s\n%s", r, buf)
		//		logger.Errorln(err)
		//	}
		//}()
		f()
	}()
}

func GoardWithChan(f func(), errChan chan error) {
	if errChan == nil {
		panic("errChan is nil")
	}
	go func() {
		defer func() {
			if r := recover(); r != nil {
				buf := make([]byte, 64<<10)
				buf = buf[:runtime.Stack(buf, false)]
				err := fmt.Errorf("GoardWithChan: recovered panic: %v\nStack trace: %s", r, buf)
				errChan <- err
			}
		}()
		f()
	}()
}

func GoardWithCallback(f func(), errCallback func(err error)) {
	if errCallback == nil {
		panic("errCallback is nil")
	}
	go func() {
		defer func() {
			if r := recover(); r != nil {
				buf := make([]byte, 64<<10)
				buf = buf[:runtime.Stack(buf, false)]
				func() {
					defer func() { recover() }()
					errCallback(fmt.Errorf("GoardWithChan: recovered panic: %v\nStack trace: %s", r, buf))
				}()
			}
		}()
		f()
	}()
}
