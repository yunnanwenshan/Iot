package logger

import (
	"fmt"
	"log"
	"os"
	"github.com/polaris1119/config"
	"strings"
	"errors"
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

var (
	logFile = "logrus.log"
	stdLogFile = "std.log"
)

var _log = New()

func init() {
	//createLogFile()
	//createStdFile();
}

//创建stdlogfile文件
func createStdFile() error {
	file, err := os.OpenFile(stdLogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Printf("open file: %s fail, err = %s", logFile, err.Error())
		//panic("open file fail")
		return errors.New("open file fail")
	} else {
		log.SetOutput(file)
		return nil
	}
}

//创建log文件
func createLogFile() error {
	file, err := getLogFile()
	if err == nil {
		_log._log.SetOutput(file)
	} else {
		fmt.Printf("Failed to log to file, using default stderr, err = %v", err)
		panic("Failed to log to file")
	}
	_log.logLevel = LevelInfo
	env, _ := config.ConfigFile.GetSection("global")
	if strings.Compare(env["env"], "debug") != 0 {
		_log.logLevel = LevelDebug
	}

	return nil
}

func getLogFile() (file *os.File, err error) {
	file, err = os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Printf("open file: %s fail, err = %s", logFile, err.Error())
		panic("open file fail")
	}

	err = nil

	return
}

//获取logger实例
func GetLoggerInstance() *logger {
	isExist := Exist(logFile)
	if isExist == false {
		createLogFile()
		//createStdFile()
	}

	return _log;
}

func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

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
	file, err := getLogFile()
	if err != nil {
		panic("open file fail")
	}
	return &logger{_log: log.New(file, "", log.Lshortfile|log.LstdFlags), logLevel: LevelDebug}
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