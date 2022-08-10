package common

import (
	"github.com/rs/zerolog"
	"os"
)

var loggerCache = make(map[string]zerolog.Logger)

func GetLoggerFor(loggerName string) zerolog.Logger {
	if logger, ok := loggerCache[loggerName]; ok {
		return logger
	} else {
		consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout}
		multi := zerolog.MultiLevelWriter(consoleWriter) //os.Stdout

		logger := zerolog.New(multi).
			With().
			Timestamp().
			Caller().
			Logger()

		loggerCache[loggerName] = logger

		return logger
	}

}
