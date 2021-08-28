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

func log(l Level, v string) {
	if l > defaultLevel {
		return
	}

	fmt.Print(v)

	if l == LevelFatal {
		os.Exit(1)
	}
}

func Fatalf(format string, v ...interface{}) {
	log(LevelFatal, fmt.Sprintf("[%s] [%s] %s", Cyan(getTime()), LightRed("FATAL"), fmt.Sprintf(format, v...)))
}

func Fatalln(v ...interface{}) {
	log(LevelFatal, fmt.Sprintf("[%s] [%s] %s", Cyan(getTime()), LightRed("FATAL"), fmt.Sprintln(v...)))
}

func Errorf(format string, v ...interface{}) {
	log(LevelError, fmt.Sprintf("[%s] [%s] %s", Cyan(getTime()), LightRed("ERROR"), fmt.Sprintf(format, v...)))
}

func Errorln(v ...interface{}) {
	log(LevelError, fmt.Sprintf("[%s] [%s] %s", Cyan(getTime()), LightRed("ERROR"), fmt.Sprintln(v...)))
}

func Infof(format string, v ...interface{}) {
	log(LevelInfo, fmt.Sprintf("[%s] [%s] %s", Cyan(getTime()), LightGreen("INFO"), fmt.Sprintf(format, v...)))
}

func Infoln(v ...interface{}) {
	log(LevelInfo, fmt.Sprintf("[%s] [%s] %s", Cyan(getTime()), LightGreen("INFO"), fmt.Sprintln(v...)))
}

func Warnf(format string, v ...interface{}) {
	log(LevelWarning, fmt.Sprintf("[%s] [%s] %s", Cyan(getTime()), Yellow("WARNING"), fmt.Sprintf(format, v...)))
}

func Warnln(v ...interface{}) {
	log(LevelWarning, fmt.Sprintf("[%s] [%s] %s", Cyan(getTime()), Yellow("WARNING"), fmt.Sprintln(v...)))
}

func Debugf(format string, v ...interface{}) {
	log(LevelDebug, fmt.Sprintf("[%s] [%s] %s", Cyan(getTime()), LightCyan("DEBUG"), fmt.Sprintf(format, v...)))
}

func Debugln(v ...interface{}) {
	log(LevelDebug, fmt.Sprintf("[%s] [%s] %s", Cyan(getTime()), LightCyan("DEBUG"), fmt.Sprintln(v...)))
}

func Verbosef(format string, v ...interface{}) {
	log(LevelVerbose, fmt.Sprintf("[%s] [%s] %s", Cyan(getTime()), LightCyan("VERBOSE"), fmt.Sprintf(format, v...)))
}

func Verboseln(format string, v ...interface{}) {
	log(LevelVerbose, fmt.Sprintf("[%s] [%s] %s", Cyan(getTime()), LightCyan("VERBOSE"), fmt.Sprintln(v...)))
}

func Successf(format string, v ...interface{}) {
	log(LevelInfo, fmt.Sprintf("[%s] [%s] %s", Cyan(getTime()), LightGreen("+"), fmt.Sprintf(format, v...)))
}

func Successln(format string, v ...interface{}) {
	log(LevelInfo, fmt.Sprintf("[%s] [%s] %s", Cyan(getTime()), LightGreen("+"), fmt.Sprintln(v...)))
}

func Failedf(format string, v ...interface{}) {
	log(LevelInfo, fmt.Sprintf("[%s] [%s] %s", Cyan(getTime()), LightGreen("-"), fmt.Sprintf(format, v...)))
}

func Failedln(v ...interface{}) {
	log(LevelInfo, fmt.Sprintf("[%s] [%s] %s", Cyan(getTime()), LightGreen("-"), fmt.Sprintln(v...)))
}

func getTime() string {
	return time.Now().Format("15:04:05")
}

func getTimestamp() string {
	return strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
}
