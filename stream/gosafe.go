package stream

import (
	"github.com/sirupsen/logrus"
)

func Recover(cleanups ...func()) {
	for _, cleanup := range cleanups {
		cleanup()
	}
	if r := recover(); r != nil {
		logrus.Error(r)
	}
}

func RunSafe(fn func()) {
	defer Recover()
	fn()
}

func GoSafe(fn func()) {
	go RunSafe(fn)
}
