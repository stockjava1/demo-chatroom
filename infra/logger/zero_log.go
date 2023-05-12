package logger

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"strings"
	"time"
)

var logger zerolog.Logger

func init() {
	//consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout}

	consoleWriter := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}
	//consoleWriter.FormatLevel = func(i interface{}) string { return strings.ToUpper(fmt.Sprintf("| %-6s|", i)) }
	//consoleWriter.FormatMessage = func(i interface{}) string { return fmt.Sprintf("***%s****", i) }
	//consoleWriter.FormatFieldName = func(i interface{}) string { return fmt.Sprintf("%s:", i) }
	//consoleWriter.FormatFieldValue = func(i interface{}) string { return strings.ToUpper(fmt.Sprintf("%s", i)) }

	// 设置 timestamp 的格式为 yyyy-MM-dd HH:mm:ss
	//consoleWriter.FormatTimestamp = func(i time.Time) string {
	//	return i.Format("2006-01-02 15:04:05")
	//}
	logger = log.Output(consoleWriter).With().Timestamp().Logger()

	//logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()
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

func main() {
	output := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}
	output.FormatLevel = func(i interface{}) string { return strings.ToUpper(fmt.Sprintf("| %-6s|", i)) }
	output.FormatMessage = func(i interface{}) string { return fmt.Sprintf("***%s****", i) }
	output.FormatFieldName = func(i interface{}) string { return fmt.Sprintf("%s:", i) }
	output.FormatFieldValue = func(i interface{}) string { return strings.ToUpper(fmt.Sprintf("%s", i)) }
	logger := log.Output(output).With().Timestamp().Logger()
	logger.Info().Str("foo", "bar").Msg("hello world")
}
