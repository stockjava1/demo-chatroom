package logger

import (
	"github.com/JabinGP/demo-chatroom/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"strings"
	"time"
)

type CustZeroLogger struct {
	logger *zerolog.Logger
	level  zerolog.Level
}

func NewLogger() *CustZeroLogger {
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
	logger := log.Output(consoleWriter).With().Caller().Timestamp().Logger()

	//logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()
	//logger = log.With().Timestamp().CallerWithSkipFrameCount(3).Stack().Logger()
	return &CustZeroLogger{
		logger: &logger,
	}
}

func NewLoggerModule(module string) *CustZeroLogger {
	logger := NewLogger()
	logger.SetModule(module)

	if config.Viper.GetString("loglevel."+module) != "" {
		logger.SetLogLevel(config.Viper.GetString("loglevel." + module))
	}

	return logger
}

func (czl *CustZeroLogger) SetModule(module string) {
	l := czl.logger.With().Str("m", module).Logger()
	czl.logger = &l
}

func (czl *CustZeroLogger) SetLogger(logger *zerolog.Logger) {
	czl.logger = logger
}

func (czl *CustZeroLogger) GetLogger() *zerolog.Logger {
	return czl.logger
}

func (czl *CustZeroLogger) SetLogLevel(level string) {
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
	//zerolog.SetGlobalLevel(l)
	czl.level = l
}

func (czl *CustZeroLogger) GetLevel() zerolog.Level {
	return czl.level
}

// Debugf implements xorm.Logger interface
func (czl *CustZeroLogger) Debug(message string, args ...interface{}) {
	if zerolog.DebugLevel >= czl.level {
		if len(args) == 0 {
			czl.logger.Debug().Msg(message)
		} else {
			czl.logger.Debug().Msgf(message, args...)
		}
	}
}

func (czl *CustZeroLogger) Info(message string, args ...interface{}) {
	if zerolog.InfoLevel >= czl.level {
		if len(args) == 0 {
			czl.logger.Info().Msg(message)
		} else {
			czl.logger.Info().Msgf(message, args...)
		}
	}
}

func (czl *CustZeroLogger) Warn(message string, args ...interface{}) {
	if zerolog.WarnLevel >= czl.level {
		if len(args) == 0 {
			czl.logger.Warn().Msg(message)
		} else {
			czl.logger.Warn().Msgf(message, args...)
		}
	}
}

func (czl *CustZeroLogger) Error(message string, args ...interface{}) {
	if zerolog.ErrorLevel >= czl.level {
		if len(args) == 0 {
			czl.logger.Error().Msg(message)
		} else {
			czl.logger.Error().Msgf(message, args...)
		}
	}
}

func (czl *CustZeroLogger) Fatal(message string, args ...interface{}) {
	if zerolog.FatalLevel >= czl.level {
		if len(args) == 0 {
			czl.logger.Fatal().Msg(message)
		} else {
			czl.logger.Fatal().Msgf(message, args...)
		}
	}
}

//func DebugJson(message string, data interface{}) {
//	logger.Debug().RawJSON(message, helpers.ServeJson(data)).Msg("")
//}

//func main() {
//	output := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}
//	output.FormatLevel = func(i interface{}) string { return strings.ToUpper(fmt.Sprintf("| %-6s|", i)) }
//	output.FormatMessage = func(i interface{}) string { return fmt.Sprintf("***%s****", i) }
//	output.FormatFieldName = func(i interface{}) string { return fmt.Sprintf("%s:", i) }
//	output.FormatFieldValue = func(i interface{}) string { return strings.ToUpper(fmt.Sprintf("%s", i)) }
//	logger := log.Output(output).With().Timestamp().Logger()
//	logger.Info().Str("foo", "bar").Msg("hello world")
//}
