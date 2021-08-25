package logger

import (
	"fmt"
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
)

func getTime() string {
	return time.Now().Format("15:04:05")
}

func getTimestamp() string {
	return strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
}

func Info(detail string) {
	fmt.Printf("[%s] [%s] %s\n", Cyan(getTime()), LightGreen("INFO"), detail)
}

func Warn(detail string) {
	fmt.Printf("[%s] [%s] %s\n", Cyan(getTime()), Yellow("WARNING"), detail)
}

func Error(detail string) {
	fmt.Printf("[%s] [%s] %s\n", Cyan(getTime()), LightRed("ERROR"), detail)
}

func Success(detail string) {
	fmt.Printf("[%s] [%s] %s\n", Cyan(getTime()), LightGreen("+"), detail)
}

func Failed(detail string) {
	fmt.Printf("[%s] [%s] %s\n", Cyan(getTime()), LightRed("-"), detail)
}
