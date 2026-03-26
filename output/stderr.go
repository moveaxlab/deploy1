package output

import (
	"strings"

	"github.com/sirupsen/logrus"
)

type ErrLogger struct{}

func (o ErrLogger) Write(p []byte) (n int, err error) {
	msg := string(p)
	if strings.Contains(msg, "error") || strings.Contains(msg, "failed") || strings.Contains(msg, "fatal") {
		logrus.Error(msg)
	} else {
		logrus.Debug(msg)
	}
	return len(p), nil
}
