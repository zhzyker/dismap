package output

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/zhzyker/dismap/internal/flag"
	"github.com/zhzyker/dismap/internal/model"
	"github.com/zhzyker/dismap/pkg/logger"
)

func Open() *os.File {
	print("openfile: ", flag.File)
	if len(flag.OutJson) != 0 {
		print(len(flag.OutJson), flag.OutJson)
		return openFile(flag.OutJson)
	} else {
		return openFile(flag.Output)
	}
}

func Write(result *model.Result, output *os.File) {
	if result.Status == "close" {
		return
	}
	if len(flag.OutJson) != 0 {
		result.Banner = hex.EncodeToString(result.BannerB)
		result.Date = time.Now()
		byteR, _ := json.Marshal(result)
		writeContent(output, string(byteR))
	} else {
		content := fmt.Sprintf("%s, %s, %s, %s, %s, %s",
			logger.GetTime(),
			result.Type,
			result.Protocol,
			logger.Clean(result.IdentifyStr),
			result.Uri,
			result.Banner)
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
