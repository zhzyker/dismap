package operate

import (
	"strconv"
	"sync"

	"github.com/zhzyker/dismap/internal/parse"
	"github.com/zhzyker/dismap/pkg/logger"
)

func FlagPing(wg *sync.WaitGroup, lock *sync.Mutex, hosts []string, TimeOut int, Thread int, np bool) []string {
	if np == true {
		return hosts
	}
	var SurviveHosts []string
	IntAllHost := 0
	IntSurHost := 0
	IntSyncHost := 0
	for _, host := range hosts {
		wg.Add(1)
		IntAllHost++
		IntSyncHost++
		go func(host string) {
			if parse.Ping(host, TimeOut) == true {
				IntSurHost++
				logger.Info("PING found alive host " + host)
				lock.Lock()
				SurviveHosts = append(SurviveHosts, host)
				lock.Unlock()
			}
			wg.Done()
		}(host)
		if IntSyncHost >= Thread {
			IntSyncHost = 0
			wg.Wait()
		}
	}
	wg.Wait()
	logger.Info(
		logger.LightGreen("There are total of ") +
			logger.White(strconv.Itoa(IntAllHost)) +
			logger.LightGreen(" hosts, and ") +
			logger.White(strconv.Itoa(IntSurHost)) +
			logger.LightGreen(" are surviving"))
	if IntSurHost <= 5 {
		logger.Warning(logger.Yellow("Too few surviving hosts"))
	}
	return SurviveHosts
}
