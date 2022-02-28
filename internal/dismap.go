package internal

import (
	"github.com/zhzyker/dismap/configs"
	"github.com/zhzyker/dismap/internal/flag"
	"github.com/zhzyker/dismap/internal/operate"
	"github.com/zhzyker/dismap/internal/output"
	"sync"

	"github.com/zhzyker/dismap/pkg/logger"
)


func which(Args map[string]interface{}, wg *sync.WaitGroup, lock *sync.Mutex) {
	op := output.Open(Args)

	address := Args["FlagNetwork"].(string)
	if address != "" {
		operate.FlagNetwork(op, wg, lock, address, Args)
		output.Close(op)
		return
	}

	uri := Args["FlagUrl"].(string)
	if uri != "" {
		operate.FlagUrl(op, uri, Args)
		output.Close(op)
		return
	}

	file := Args["FlagFile"].(string)
	if file != "" {
		operate.FlagFile(op, wg, lock, file, Args)
		output.Close(op)
		return
	}

}

func DisMap() {
	configs.Banner()
	Args := flag.Flags()
	wg := &sync.WaitGroup{}
	lock := &sync.Mutex{}

	information()
	which(Args, wg, lock)
	logger.Info("Identification completed and ended")
}
