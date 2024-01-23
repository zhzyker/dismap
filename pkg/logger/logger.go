package logger

import (
	"fmt"
	"path/filepath"
	"runtime"
	"time"
)

// logMsg 函数用于记录日志信息
func logMsg(msg string, level int) {
	// 判断日志级别, 小于当前日志级别则不输出日志, 配合 SetLogLevel 设定日志级别
	if level < logLevel {
		return
	}
	// 获取日志发生的位置信息
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	}
	// 根据日志级别获取对应的日志级别名称
	var levelName string
	switch level {
	case DEBUG:
		levelName = "DBG"
	case INFO:
		levelName = "INF"
	case WARN:
		levelName = "WAR"
	case ERROR:
		levelName = "ERR"
	default:
		levelName = "???"
	}
	// 格式化日志信息
	t := time.Now().Format("2006/01/02 15:04:05")
	content := fmt.Sprintf("%s [%s] [%s:%d] %s\n", t, levelName, filepath.Base(file), line, msg)
	// 直接显示日志信息
	fmt.Print(content)
	// writeLogToFile 交给写入模块把日志信息写入文件
	// writeLogToFile(content)
}

func INF(msg string) {
	logMsg(msg, INFO)
}

func WAR(msg string) {
	logMsg(msg, WARN)
}

func ERR(msg string) {
	logMsg(msg, ERROR)
}

func DBG(msg string) {
	logMsg(msg, DEBUG)
}
