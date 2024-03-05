package output

import "github.com/sirupsen/logrus"

type ErrLogger struct{}

func (o ErrLogger) Write(p []byte) (n int, err error) {
	logrus.Error(string(p))
	return len(p), nil
}
