package operate

import (
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/zhzyker/dismap/internal/flag"
	"github.com/zhzyker/dismap/internal/output"
	"github.com/zhzyker/dismap/internal/parse"
	"github.com/zhzyker/dismap/internal/protocol"
	"github.com/zhzyker/dismap/pkg/logger"
)

func FlagNetwork(op *os.File, wg *sync.WaitGroup, lock *sync.Mutex, address string) {
	timeout := flag.Timeout
	thread := flag.Thread
	np := flag.NoIcmp
	flagPort := flag.Port
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
			go func(host string, port int) {
				resTls := protocol.DiscoverTls(host, port)
				if resTls.Status == "open" {
					intAll++
					parse.VerboseParse(resTls)
					output.Write(resTls, op)
					if strings.Contains(resTls.Uri, "://") {
						intIde++
					}
				}
				wg.Done()
			}(host, port)

			go func(host string, port int) {
				resTcp := protocol.DiscoverTcp(host, port)
				if resTcp.Status == "open" {
					intAll++
					parse.VerboseParse(resTcp)
					output.Write(resTcp, op)
					if strings.Contains(resTcp.Uri, "://") {
						intIde++
					}
				}
				wg.Done()
			}(host, port)

			go func(host string, port int) {
				resUdp := protocol.DiscoverUdp(host, port)
				if resUdp.Status == "open" {
					intAll++
					parse.VerboseParse(resUdp)
					output.Write(resUdp, op)
					if strings.Contains(resUdp.Uri, "://") {
						intIde++
					}
				}
				wg.Done()
			}(host, port)
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
