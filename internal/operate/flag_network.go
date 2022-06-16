package operate

import (
	"github.com/zhzyker/dismap/internal/output"
	"github.com/zhzyker/dismap/internal/parse"
	"github.com/zhzyker/dismap/internal/protocol"
	"github.com/zhzyker/dismap/pkg/logger"
	"os"
	"strconv"
	"strings"
	"sync"
)

func FlagNetwork(op *os.File, wg *sync.WaitGroup, lock *sync.Mutex, address string, Args map[string]interface{}) {
	timeout := Args["FlagTimeout"].(int)
	thread := Args["FlagThread"].(int)
	np := Args["FlagNoIcmp"].(bool)
	flagPort := Args["FlagPort"].(string)
	ports := parse.PortParse(flagPort)

	logger.Info("Start to detect host from " + address)
	hosts, err := parse.NetworkParse(address)
	if logger.DebugError(err) {
		return
	}

	actualHosts := FlagPing(wg, lock, hosts, timeout, thread, np)

	logger.Info("Start to identify the targets")
	intSyncThread := 0
	intAll := 0
	intIde := 0
	for _, host := range actualHosts {
		for _, port := range ports {
			wg.Add(3)
			intSyncThread++
			go func(host string, port int, Args map[string]interface{}) {
				resTls := protocol.DiscoverTls(host, port, Args)
				if resTls["status"].(string) == "open" {
					intAll++
					parse.VerboseParse(resTls)
					output.Write(resTls, op)
					if strings.Contains(resTls["uri"].(string), "://") {
						intIde++
					}
				}
				wg.Done()
			}(host, port, Args)

			go func(host string, port int, Args map[string]interface{}) {
				resTcp := protocol.DiscoverTcp(host, port, Args)
				if resTcp["status"].(string) == "open" {
					intAll++
					parse.VerboseParse(resTcp)
					output.Write(resTcp, op)
					if strings.Contains(resTcp["uri"].(string), "://") {
						intIde++
					}
				}
				wg.Done()
			}(host, port, Args)

			go func(host string, port int, Args map[string]interface{}) {
				resUdp := protocol.DiscoverUdp(host, port, Args)
				if resUdp["status"].(string) == "open" {
					intAll++
					parse.VerboseParse(resUdp)
					output.Write(resUdp, op)
					if strings.Contains(resUdp["uri"].(string), "://") {
						intIde++
					}
				}
				wg.Done()
			}(host, port, Args)
			if intSyncThread >= thread {
				intSyncThread = 0
				wg.Wait()
			}
		}
	}
	wg.Wait()
	logger.Info(logger.LightGreen("A total of ") +
		logger.White(strconv.Itoa(intAll)) +
		logger.LightGreen(" targets, the rule base hits ") +
		logger.White(strconv.Itoa(intIde)) +
		logger.LightGreen(" targets"))
}
