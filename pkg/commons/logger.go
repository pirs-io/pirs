package commons

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

var loggerCache = make(map[string]zerolog.Logger)

func InitLogger(logPath string, maxBackups int, maxSize int, maxAge int) {
	logFile := &lumberjack.Logger{
		Filename:   logPath,
		MaxBackups: maxBackups, // recommended 10
		MaxSize:    maxSize,    //recommended 3
		MaxAge:     maxAge,     //recommended 30
		Compress:   true,
	}

	log.Logger = log.Output(zerolog.MultiLevelWriter(zerolog.ConsoleWriter{Out: os.Stdout}, logFile)).With().Timestamp().Caller().Logger()
}

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

func CheckAndLog(err error, log zerolog.Logger) error {
	if err != nil {
		log.Err(err).Msg(err.Error())
		return err
	}
	return nil
}
