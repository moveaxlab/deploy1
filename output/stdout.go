package output

import "github.com/sirupsen/logrus"

type OutLogger struct{}

func (o OutLogger) Write(p []byte) (n int, err error) {
	logrus.Debug(string(p))
	return len(p), nil
}
