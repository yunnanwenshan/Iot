package log

import (
	"fmt"
	"log"
	"os"
)

type (
	Level int
)

const (
	LevelFatal = iota
	LevelError
	LevelWarning
	LevelInfo
	LevelDebug
)

var _log *logger = New()

func Fatal(s string) {
	_log.Output(LevelFatal, s)
	os.Exit(1)
}

func Fatalf(format string, v ...interface{}) {
	_log.Output(LevelFatal, fmt.Sprintf(format, v...))
	os.Exit(1)
}

func Error(s string) {
	_log.Output(LevelError, s)
}

func Errorf(format string, v ...interface{}) {
	_log.Output(LevelError, fmt.Sprintf(format, v...))
}

func Warn(s string) {
	_log.Output(LevelWarning, s)
}

func Warnf(format string, v ...interface{}) {
	_log.Output(LevelWarning, fmt.Sprintf(format, v...))
}

func Info(s string) {
	_log.Output(LevelInfo, s)
}

func Infof(format string, v ...interface{}) {
	_log.Output(LevelInfo, fmt.Sprintf(format, v...))
}

func Debug(s string) {
	_log.Output(LevelDebug, s)
}

func Debugf(format string, v ...interface{}) {
	_log.Output(LevelDebug, fmt.Sprintf(format, v...))
}

func SetLogLevel(level Level) {
	_log.SetLogLevel(level)
}

type logger struct {
	_log *log.Logger
	//小于等于该级别的level才会被记录
	logLevel Level
}

//NewLogger 实例化，供自定义
func NewLogger() *logger {
	return &logger{_log: log.New(os.Stderr, "", log.Lshortfile|log.LstdFlags), logLevel: LevelDebug}
}

//New 实例化，供外部直接调用 log.XXXX
func New() *logger {
	return &logger{_log: log.New(os.Stderr, "", log.Lshortfile|log.LstdFlags), logLevel: LevelDebug}
}

func (l *logger) Output(level Level, s string) error {
	if l.logLevel < level {
		return nil
	}
	formatStr := "[UNKNOWN] %s"
	switch level {
	case LevelFatal:
		formatStr = "\033[35m[FATAL]\033[0m %s"
	case LevelError:
		formatStr = "\033[31m[ERROR]\033[0m %s"
	case LevelWarning:
		formatStr = "\033[33m[WARN]\033[0m %s"
	case LevelInfo:
		formatStr = "\033[32m[INFO]\033[0m %s"
	case LevelDebug:
		formatStr = "\033[36m[DEBUG]\033[0m %s"
	}
	s = fmt.Sprintf(formatStr, s)
	return l._log.Output(3, s)
}

func (l *logger) Fatal(s string) {
	l.Output(LevelFatal, s)
	os.Exit(1)
}

func (l *logger) Fatalf(format string, v ...interface{}) {
	l.Output(LevelFatal, fmt.Sprintf(format, v...))
	os.Exit(1)
}

func (l *logger) Error(s string) {
	l.Output(LevelError, s)
}

func (l *logger) Errorf(format string, v ...interface{}) {
	l.Output(LevelError, fmt.Sprintf(format, v...))
}

func (l *logger) Warn(s string) {
	l.Output(LevelWarning, s)
}

func (l *logger) Warnf(format string, v ...interface{}) {
	l.Output(LevelWarning, fmt.Sprintf(format, v...))
}

func (l *logger) Info(s string) {
	l.Output(LevelInfo, s)
}

func (l *logger) Infof(format string, v ...interface{}) {
	l.Output(LevelInfo, fmt.Sprintf(format, v...))
}

func (l *logger) Debug(s string) {
	l.Output(LevelDebug, s)
}

func (l *logger) Debugf(format string, v ...interface{}) {
	l.Output(LevelDebug, fmt.Sprintf(format, v...))
}

func (l *logger) SetLogLevel(level Level) {
	l.logLevel = level
}
