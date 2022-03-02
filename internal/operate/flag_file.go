package operate

import (
	"bufio"
	"github.com/zhzyker/dismap/internal/parse"
	"github.com/zhzyker/dismap/pkg/logger"
	"io"
	"net/url"
	"os"
	"strings"
	"sync"
)

func FlagFile(op *os.File, wg *sync.WaitGroup, lock *sync.Mutex, file string, Args map[string]interface{}) {
	thread := Args["FlagThread"].(int)
	f, err := os.Open(file)
	if err != nil {
		logger.Error("There is no " + logger.LightRed(f) + " file or the directory does not exist")
	}

	logger.Info(logger.LightGreen("Batch scan the targets in " + logger.Yellow(file) + logger.LightGreen(", priority network segment")))
	buf := bufio.NewReader(f)

	intSyncThread := 0
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		if logger.DebugError(err) || err == io.EOF {
			break
		}
		if line == "" {
			continue
		}

		if parse.NetJudgeParse(line) {
			FlagNetwork(op, wg, lock, line, Args)
			continue
		}
		_, err = url.Parse(line)
		if logger.DebugError(err) {
			logger.Error(logger.Red("Unable to parse: " + line))
			continue
		} else {
			wg.Add(1)
			intSyncThread++
			go func(line string, Args map[string]interface{}) {
				lock.Lock()
				FlagUrl(op, line, Args)
				lock.Unlock()
				wg.Done()
			}(line, Args)
			if intSyncThread >= thread {
				intSyncThread = 0
				wg.Wait()
			}
			continue
		}
	}
	wg.Wait()
}