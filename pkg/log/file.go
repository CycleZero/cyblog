package log

import (
	"io"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

var globalLogWriter io.Writer

func NewFileWriter(logPath string) *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename: logPath,
		MaxSize:  10, // megabytes
	}
}

func GetLogPath(logDir string) string {
	if logDir == "" {
		return ""
	}
	p := logDir + "/" + time.Now().Format("2006-01-02") + "/" + time.Now().Format("2006-01-02-15-04-05") + ".log"
	return p
}

func GetLogWriter() io.Writer {
	if globalLogWriter == nil {
		panic("log file writer 未初始化")
	}
	return globalLogWriter
}

func InitGlobalLogWriter(consoleWriter io.Writer, fileWriter io.Writer) {
	if consoleWriter == nil {
		panic("console writer 未初始化")
	}
	if fileWriter == nil {
		globalLogWriter = consoleWriter
	} else {
		globalLogWriter = io.MultiWriter(consoleWriter, fileWriter)
	}
}
