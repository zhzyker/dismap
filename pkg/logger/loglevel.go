package logger

const (
	DEBUG = iota
	INFO
	WARN
	ERROR
)

var logLevel = INFO // 默认日志级别为 INFO

// SetLogLevel 用于设置日志级别
// 参数 level 为要设置的日志级别, 取值范围为常量 logger.DEBUG logger.INFO logger.WARN logger.ERROR
func SetLogLevel(level int) {
	logLevel = level
}
