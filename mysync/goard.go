package mysync

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"runtime"
)

var logger *logrus.Entry = logrus.NewEntry(logrus.StandardLogger())

// ChangeLogger change logger to newLogger, use nil to prevent Goard from logging
func ChangeLogger(newLogger *logrus.Entry) {
	logger = newLogger
}

func Goard(ctx context.Context, f func(context.Context)) {
	go func() {
		defer func() {
			recover()
		}()
		f(ctx)
	}()
}

func GoardWithLog(ctx context.Context, f func(context.Context)) {
	go func() {
		defer func() {
			if r := recover(); logger != nil && r != nil {
				buf := make([]byte, 64<<10)
				buf = buf[:runtime.Stack(buf, false)]
				err := fmt.Errorf("errgroup: panic recovered: %s\n%s", r, buf)
				logger.Errorln(err)
			}
		}()
		f(ctx)
	}()
}

func GoardWithLogger(ctx context.Context, f func(context.Context), logger *logrus.Logger) {
	if logger == nil {
		panic("logger is nil")
	}
	go func() {
		defer func() {
			if r := recover(); logger != nil && r != nil {
				buf := make([]byte, 64<<10)
				buf = buf[:runtime.Stack(buf, false)]
				err := fmt.Errorf("errgroup: panic recovered: %s\n%s", r, buf)
				logger.Errorln(err)
			}
		}()
		f(ctx)
	}()
}

func GoardWithChan(ctx context.Context, f func(context.Context), errChan chan error) {
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
		f(ctx)
	}()
}

func GoardWithCallback(ctx context.Context, f func(context.Context), errCallback func(err error)) {
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
		f(ctx)
	}()
}
