package logger

import (
	"github.com/sirupsen/logrus"
)

func SetIdentifierField(identifier string) logrus.Entry {
	customLog := logrus.WithFields(logrus.Fields{
		"identifier": identifier,
	})
	return *customLog
}