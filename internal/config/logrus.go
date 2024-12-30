package config

import (
	"strconv"

	"github.com/sirupsen/logrus"
)

func NewLogger(conf *Config) *logrus.Logger {
	log := logrus.New()
	logLevel, _ := strconv.Atoi(conf.Logger.Level)
	log.SetLevel(logrus.Level(logLevel))
	log.SetFormatter(&logrus.JSONFormatter{})

	return log
}
