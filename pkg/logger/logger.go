package logger

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gookit/color"
)

var (
	Cyan       = color.Cyan.Render
	Yellow     = color.Yellow.Render
	White      = color.White.Render
	Blue       = color.Blue.Render
	Purple     = color.Style{color.Magenta, color.OpBold}.Render
	LightRed   = color.Style{color.Red, color.OpBold}.Render
	LightGreen = color.Style{color.Green, color.OpBold}.Render
	LightWhite = color.Style{color.White, color.OpBold}.Render
	LightCyan  = color.Style{color.Cyan, color.OpBold}.Render
)

var (
	defaultLevel = LevelInfo
)

func SetLevel(l Level) {
	defaultLevel = l
}

func log(l Level, detail string) {
	if l > defaultLevel {
		return
	}

	fmt.Println(detail)

	if l == LevelFatal {
		os.Exit(1)
	}
}

func Fatal(detail string) {
	log(LevelFatal, fmt.Sprintf("[%s] [%s] %s", Cyan(getTime()), LightRed("FATAL"), detail))
}

func Error(detail string) {
	log(LevelError, fmt.Sprintf("[%s] [%s] %s", Cyan(getTime()), LightRed("ERROR"), detail))
}

func Info(detail string) {
	log(LevelInfo, fmt.Sprintf("[%s] [%s] %s", Cyan(getTime()), LightGreen("INFO"), detail))
}

func Warn(detail string) {
	log(LevelWarning, fmt.Sprintf("[%s] [%s] %s", Cyan(getTime()), Yellow("WARNING"), detail))
}

func Debug(detail string) {
	log(LevelDebug, fmt.Sprintf("[%s] [%s] %s", Cyan(getTime()), LightCyan("DEBUG"), detail))
}

func Verbose(detail string) {
	log(LevelVerbose, fmt.Sprintf("[%s] [%s] %s", Cyan(getTime()), LightCyan("VERBOSE"), detail))
}

func Success(detail string) {
	log(LevelInfo, fmt.Sprintf("[%s] [%s] %s", Cyan(getTime()), LightGreen("+"), detail))
}

func Failed(detail string) {
	log(LevelInfo, fmt.Sprintf("[%s] [%s] %s", Cyan(getTime()), LightGreen("-"), detail))
}

func getTime() string {
	return time.Now().Format("15:04:05")
}

func getTimestamp() string {
	return strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
}
