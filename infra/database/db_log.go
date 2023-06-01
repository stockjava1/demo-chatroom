package database

import (
	"fmt"
	"github.com/JabinGP/demo-chatroom/infra/logger"
	"github.com/rs/zerolog"
	"xorm.io/core"
)

func printStrings(v ...interface{}) string {
	// 创建一个字符串切片，用于存储每个参数的字符串表示
	s := make([]string, len(v))
	// 遍历每个参数
	for i, val := range v {
		// 根据参数的类型，使用 fmt.Sprintf 函数并指定相应的格式化动词来获取参数的字符串表示，并存入切片
		switch val.(type) {
		case string:
			s[i] = fmt.Sprintf("string: %s", val)
		case int:
			s[i] = fmt.Sprintf("int: %d", val)
		case bool:
			s[i] = fmt.Sprintf("bool: %t", val)
		case []int:
			s[i] = fmt.Sprintf("[]int: %v", val)
		case struct{ Name string }:
			s[i] = fmt.Sprintf("struct: %+v", val)
		default:
			s[i] = fmt.Sprintf("unknown: %v", val)
		}
	}
	// 打印切片中的所有字符串
	return fmt.Sprintf("%v\n", s)
}

// ZerologLogger implements xorm.Logger interface with zerolog
type DbLogger struct {
	logger logger.CustZeroLogger
	level  core.LogLevel
}

// NewZerologLogger creates a new ZerologLogger instance
func NewDbLogger(logger *logger.CustZeroLogger, l core.LogLevel) *DbLogger {

	dbLogger := &DbLogger{
		logger: *logger,
		level:  l,
	}
	dbLogger.logger.SetLevel(level(l))

	return dbLogger
}

func level(l core.LogLevel) zerolog.Level {
	var nl zerolog.Level
	switch l {
	case core.LOG_ERR:
		nl = zerolog.ErrorLevel
	case core.LOG_WARNING:
		nl = zerolog.WarnLevel
	case core.LOG_INFO:
		nl = zerolog.InfoLevel
	case core.LOG_DEBUG:
		nl = zerolog.DebugLevel
	default:
		nl = zerolog.InfoLevel
	}

	return nl
}

// Debugf implements xorm.Logger interface
func (zl *DbLogger) Debugf(format string, v ...interface{}) {
	if zl.level <= core.LOG_DEBUG {
		zl.logger.Debug().Msgf(format, v...)
	}
}

// Infof implements xorm.Logger interface
func (zl *DbLogger) Infof(format string, v ...interface{}) {
	if zl.level <= core.LOG_INFO {
		zl.logger.Info().Msgf(format, v...)
	}
}

// Warnf implements xorm.Logger interface
func (zl *DbLogger) Warnf(format string, v ...interface{}) {
	if zl.level <= core.LOG_WARNING {
		zl.logger.Warn().Msgf(format, v...)
	}
}

// Errorf implements xorm.Logger interface
func (zl *DbLogger) Errorf(format string, v ...interface{}) {
	if zl.level <= core.LOG_ERR {
		zl.logger.Error().Msgf(format, v...)
	}
}

// Debugf implements xorm.Logger interface
func (zl *DbLogger) Debug(v ...interface{}) {
	if zl.level <= core.LOG_DEBUG {
		zl.logger.Debug().Msg(printStrings(v))
	}
}

// Infof implements xorm.Logger interface
func (zl *DbLogger) Info(v ...interface{}) {
	if zl.level <= core.LOG_INFO {
		zl.logger.Info().Msg(printStrings(v))
	}
}

// Warnf implements xorm.Logger interface
func (zl *DbLogger) Warn(v ...interface{}) {
	if zl.level <= core.LOG_WARNING {
		zl.logger.Warn().Msg(printStrings(v))
	}
}

// Errorf implements xorm.Logger interface
func (zl *DbLogger) Error(v ...interface{}) {
	if zl.level <= core.LOG_ERR {
		zl.logger.Error().Msg(printStrings(v))
	}
}

// Level implements xorm.Logger interface
func (zl *DbLogger) Level() core.LogLevel {
	return zl.level
}

// SetLevel implements xorm.Logger interface
func (zl *DbLogger) SetLevel(l core.LogLevel) {
	zl.level = l
	zl.logger.SetLevel(level(l))
}

// ShowSQL implements xorm.Logger interface
func (zl *DbLogger) ShowSQL(show ...bool) {}

// IsShowSQL implements xorm.Logger interface
func (zl *DbLogger) IsShowSQL() bool {
	return true
}
