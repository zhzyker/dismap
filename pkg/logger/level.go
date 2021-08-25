package logger

type Level int

const (
	LevelFatal Level = iota
	LevelError
	LevelInfo
	LevelWarning
	LevelDebug
	LevelVerbose
)
