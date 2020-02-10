package internal // import "golang.handcraftedbits.com/ezif/internal"

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

//
// Public variables
//

var Log *logrus.Logger

//
// Private functions
//

func init() {
	var formatter logrus.Formatter
	var level logrus.Level

	formatter, level = newLogConfig(os.Getenv("EZIF_LOG_FORMAT"), os.Getenv("EZIF_LOG_LEVEL"))

	Log = logrus.New()

	Log.Formatter = formatter
	Log.Level = level
}

func newLogConfig(formatter, level string) (logrus.Formatter, logrus.Level) {
	var err error
	var parsedFormatter logrus.Formatter
	var parsedLevel logrus.Level

	parsedLevel, err = logrus.ParseLevel(level)

	if err != nil {
		parsedLevel = logrus.InfoLevel
	}

	switch strings.ToLower(strings.TrimSpace(formatter)) {
	case "json":
		parsedFormatter = &logrus.JSONFormatter{}
	default:
		parsedFormatter = &logrus.TextFormatter{
			EnvironmentOverrideColors: true,
		}
	}

	return parsedFormatter, parsedLevel
}
