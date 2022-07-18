package output

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/zhzyker/dismap/internal/flag"
	"github.com/zhzyker/dismap/pkg/logger"
	"os"
	"time"
)

func Open(Args map[string]interface{}) *os.File {
	if len(Args["FlagOutJson"].(string)) != 0 {
		return openFile(Args["FlagOutJson"].(string))
	} else {
		return openFile(Args["FlagOutput"].(string))
	}
}

func Write(result map[string]interface{}, output *os.File) {
	if result["status"].(string) == "close" {
		return
	}
	if len(flag.OutJson) != 0 {
		result["banner.byte"] = hex.EncodeToString(result["banner.byte"].([]byte))
		result["date"] = time.Now().Unix()
		byteR, _ := json.Marshal(result)
		writeContent(output, string(byteR))
	} else {
		content := fmt.Sprintf("%s, %s, %d, %s, %s, %s, %s, %s",
			logger.GetTime(),
			result["host"],
			result["port"],
			result["type"],
			result["protocol"],
			logger.Clean(result["identify.string"].(string)),
			result["uri"],
			result["banner.string"])
		writeContent(output, content)
		writeContent(output, content)
	}
}

func Close(file *os.File) {
	err := file.Close()
	if logger.DebugError(err) {
		logger.Error(fmt.Sprintf("Close file %s exception", logger.Red(file.Name())))
	} else {
		logger.Info("The identification results are saved in " + logger.Yellow(file.Name()))
	}
}

func openFile(name string) *os.File {
	osFile, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if logger.DebugError(err) {
		logger.Error(fmt.Sprintf("Failed to open file %s", logger.Red(name)))
	}
	return osFile
}

func writeContent(file *os.File, content string) {
	_, err := file.Write([]byte(content + "\n"))
	if logger.DebugError(err) {
		logger.Error(fmt.Sprintf("Write failed: %s", logger.Red(content)))
	}
}
