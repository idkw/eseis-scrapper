package logger

import (
	"github.com/sirupsen/logrus"
)

func init() {
	// Only log the warning severity or above.
	logrus.SetLevel(logrus.DebugLevel)
}
