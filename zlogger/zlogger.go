package zlogger

import (
	"errors"
	"strings"
)

type zLogLevel uint16

const (
	// 为不同的level赋值 UNKNOWN=0
	UNKNOWN zLogLevel = iota
	DEBUG
	TRACE
	INFO
	WARN
	ERROR
	FATAL
)

type ZLoggerInterface interface {
	Debug(msg string, arg ...interface{})
	Info(msg string, arg ...interface{})
	Warn(msg string, arg ...interface{})
	Error(msg string, arg ...interface{})
	Fatal(msg string, arg ...interface{})
}

func parseLogLevel(s string) (zLogLevel, error) {
	s = strings.ToLower(s)
	switch s {
	case "debug":
		return DEBUG, nil
	case "trace":
		return TRACE, nil
	case "info":
		return INFO, nil
	case "warn":
		return WARN, nil
	case "error":
		return ERROR, nil
	case "fatal":
		return FATAL, nil
	default:
		err := errors.New("logger level string is valid")
		return UNKNOWN, err
	}
}

func parseLogLevelToStr(level zLogLevel) (string, error) {
	switch level {
	case DEBUG:
		return "DEBUG", nil
	case TRACE:
		return "TRACE", nil
	case INFO:
		return "INFO", nil
	case WARN:
		return "WARN", nil
	case ERROR:
		return "ERROR", nil
	case FATAL:
		return "FATAL", nil
	default:
		err := errors.New("logger level is valid")
		return "UNKNOWN", err
	}
}

func getLogString(lv zLogLevel) string {
	switch lv {
	case DEBUG:
		return Cyan("DEBUG")
	case TRACE:
		return Purple("TRACE")
	case INFO:
		return Green("INFO")
	case WARN:
		return Yellow("WARN")
	case ERROR:
		return Red("ERROR")
	case FATAL:
		return Blue("FATAL")
	default:
		return White("UNKNOWN")
	}
}
