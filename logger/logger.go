package logger

import (
	"fmt"
	"github.com/lemontree2015/skynet/config"
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

var Logger *logrus.Logger

var loggerLevel map[string]logrus.Level

func init() {
	loggerLevel = make(map[string]logrus.Level)
	loggerLevel["panic"] = logrus.PanicLevel
	loggerLevel["fatal"] = logrus.FatalLevel
	loggerLevel["error"] = logrus.ErrorLevel
	loggerLevel["warn"] = logrus.WarnLevel
	loggerLevel["info"] = logrus.InfoLevel
	loggerLevel["debug"] = logrus.DebugLevel
	loggerLevel["trace"] = logrus.TraceLevel

	Logger = logrus.New()
	//设置输出样式，自带的只有两种样式logrus.JSONFormatter{}和logrus.TextFormatter{}
	Logger.SetFormatter(&logrus.JSONFormatter{})
	Logger.SetOutput(os.Stdout)

	logPath, err := config.DefaultString("log.path")
	if err != nil {
		panic(fmt.Errorf("log.path is not set"))
	}
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	//file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	writers := []io.Writer{file, os.Stdout}
	//同时写文件和屏幕
	fileAndStdoutWriter := io.MultiWriter(writers...)
	if err == nil {
		Logger.SetOutput(fileAndStdoutWriter)
	} else {
		Logger.Info("failed to log to file.")
	}
	cfgLogLevel, err := config.DefaultString("log.level")
	if err != nil {
		cfgLogLevel = "info"
	}
	//设置最低loglevel
	if level, ok := loggerLevel[cfgLogLevel]; ok {
		Logger.SetLevel(level)
	} else {
		Logger.SetLevel(logrus.InfoLevel)
	}
}
