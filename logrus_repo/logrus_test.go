package logrus_repo

import (
	"io/ioutil"
	"testing"
	"time"

	nested "github.com/antonfisher/nested-logrus-formatter"
	logredis "github.com/rogierlommers/logrus-redis-hook"
	"github.com/sirupsen/logrus"
)

type AppHook struct {
	AppName string
}

func (h *AppHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *AppHook) Fire(entry *logrus.Entry) error {
	entry.Data["app"] = h.AppName
	return nil
}

func TestReportCaller(t *testing.T) {
	logrus.SetReportCaller(true)
	logrus.Infof("info msg")
}

func TestWithFields(t *testing.T) {
	logrus.WithFields(logrus.Fields{
		"name": "dj",
		"age":  18,
	}).Info("info msg")
}

func TestCustomLogger(t *testing.T) {
	log := logrus.New()

	log.SetLevel(logrus.InfoLevel)
	log.SetFormatter(&logrus.TextFormatter{})

	log.Info("info msg")
	log.Debug("debug msg")

	log.SetFormatter(&logrus.JSONFormatter{})
	log.Info("info msg json format")

	log.SetFormatter(&nested.Formatter{
		TimestampFormat: time.RFC3339,
		HideKeys:        true,
		FieldsOrder:     []string{"component", "category"},
	})
	log.Info("info msg")
}

func TestHook(t *testing.T) {
	log := logrus.New()
	h := &AppHook{
		AppName: "awesome-web",
	}
	log.AddHook(h)
	log.Infof("Add hook")
}

func TestRedisHook(t *testing.T) {
	log := logrus.New()
	hookConfig := logredis.HookConfig{
		Host:     "localhost",
		Key:      "mykey",
		Format:   "v0",
		App:      "awesome",
		Hostname: "localhost",
		TTL:      3600,
		Port:     6379,
	}

	hook, err := logredis.NewHook(hookConfig)
	if err == nil {
		log.AddHook(hook)
	} else {
		log.Errorf("logredis error: %q", err)
	}

	log.Info("just some info logging...")

	log.WithFields(logrus.Fields{
		"animal": "walrus",
		"foo":    "bar",
		"this":   "that",
	}).Info("additional fields are being logged as well")

	log.SetOutput(ioutil.Discard)
	log.Info("This will only be sent to Redis")
}
