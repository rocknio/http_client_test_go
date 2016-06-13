package main

import (
	"fmt"
	"os"

	"github.com/op/go-logging"
)

// Logger 日志模块
var Logger = logging.MustGetLogger("http_client_test")

var format = logging.MustStringFormatter(
	`%{time:2006-01-02 15:04:05.999999} %{shortfile} - %{shortfunc} ▶ %{level:.4s} %{message}`,
)

// InitLogger 初始化日志模块
func InitLogger(logfilename string) (err error) {
	// 打开日志文件
	logfile, logerr := os.OpenFile(logfilename, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if logerr != nil {
		fmt.Println("Fail to Create/Open", logfilename, "Http Client Test start Failed! ERR = ", logerr)
		return logerr
	}

	backend := logging.NewLogBackend(logfile, "", 0)
	backendFormater := logging.NewBackendFormatter(backend, format)

	backend1 := logging.NewLogBackend(os.Stderr, "", 0)
	backend1Formatter := logging.NewBackendFormatter(backend1, format)

	logging.SetBackend(backendFormater, backend1Formatter)

	return nil
}

// SetLogLevel 设置日志级别
// DEBUG
// NOTICE
// WARNNING
// ERROR
// FATAL
func SetLogLevel(loglevel string) {
	switch loglevel {
	case "DEBUG":
		logging.SetLevel(logging.DEBUG, loglevel)
	case "NOTICE":
		logging.SetLevel(logging.NOTICE, loglevel)
	case "WARNNING":
		logging.SetLevel(logging.WARNING, loglevel)
	case "ERROR":
		logging.SetLevel(logging.ERROR, loglevel)
	case "FATAL":
		logging.SetLevel(logging.CRITICAL, loglevel)
	default:
	case "INFO":
		logging.SetLevel(logging.INFO, loglevel)
	}
}
