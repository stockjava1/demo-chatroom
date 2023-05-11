package logger

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"strings"
)

var logger zerolog.Logger

func init() {
	logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()
	//logger = log.With().Timestamp().CallerWithSkipFrameCount(3).Stack().Logger()
}

func Logger() zerolog.Logger {
	return logger
}

func SetLogLevel(level string) {
	var l zerolog.Level

	switch strings.ToLower(level) {
	case "error":
		l = zerolog.ErrorLevel
	case "warn":
		l = zerolog.WarnLevel
	case "info":
		l = zerolog.InfoLevel
	case "debug":
		l = zerolog.DebugLevel
	default:
		l = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(l)

}

func Info(message string, args ...interface{}) {
	logger.Info().Msgf(message, args...)
}

//func DebugJson(message string, data interface{}) {
//	logger.Debug().RawJSON(message, helpers.ServeJson(data)).Msg("")
//}

func Debug(message string, args ...interface{}) {
	logger.Debug().Msgf(message, args...)
}

func Warn(message string, args ...interface{}) {
	logger.Warn().Msgf(message, args...)
}

func Error(message string, args ...interface{}) {
	logger.Error().Msgf(message, args...)
}

func Fatal(message string, args ...interface{}) {
	logger.Fatal().Msgf(message, args)
	os.Exit(1)
}

func Log(message string, args ...interface{}) {
	if len(args) == 0 {
		logger.Info().Msg(message)
	} else {
		logger.Info().Msgf(message, args...)
	}
}
