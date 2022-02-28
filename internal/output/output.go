package output

import (
	"fmt"
	"github.com/zhzyker/dismap/pkg/logger"
	"os"
)

func Open(Args map[string]interface{}) *os.File  {
	o := Args["FlagOutput"].(string)
	op, err :=  os.OpenFile(o, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if logger.DebugError(err) {
		logger.Error(fmt.Sprintf("Failed to open file %s", logger.Red(o)))
	}
	return op
}

func Write(result map[string]interface{}, output *os.File) {
	if result["status"].(string) == "close" {
		return
	}
	JsonOutput(result, "save")
	content := fmt.Sprintf("%s, %s, %s, %s, %s, %s",
		logger.GetTime(),
		result["type"],
		result["protocol"],
		logger.Clean(result["identify.string"].(string)),
		result["uri"],
		result["banner.string"])
	var text = []byte(content + "\n")
	_, err := output.Write(text)
	if logger.DebugError(err) {
		logger.Error(fmt.Sprintf("Target %s write failed", logger.Red(result["uri"])))
	}
}

func Close(file *os.File) {
	JsonOutput(nil, "write")
	err := file.Close()
	if logger.DebugError(err) {
		logger.Error(fmt.Sprintf("Close file %s exception", logger.Red(file.Name())))
	} else {
		logger.Info("The identification results are saved in " + logger.Yellow(file.Name()))
	}
}
