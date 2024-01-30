package internal

import (
	"net/http"
	"os"
	"sync"

	"github.com/ivanpirog/coloredcobra"
	"github.com/spf13/cobra"
	"github.com/zhzyker/dismap/configs"
	"github.com/zhzyker/dismap/internal/flag"
	"github.com/zhzyker/dismap/internal/operate"
	"github.com/zhzyker/dismap/internal/output"

	"github.com/zhzyker/dismap/pkg/logger"
)

func which(wg *sync.WaitGroup, lock *sync.Mutex) {
	op := output.Open()

	address := flag.NetWork
	if address != "" {
		operate.FlagNetwork(op, wg, lock, address)
		output.Close(op)
		return
	}

	uri := flag.InUrl
	if uri != "" {
		operate.FlagUrl(op, uri)
		output.Close(op)
		return
	}

	file := flag.File
	if file != "" {
		operate.FlagFile(op, wg, lock, file)
		output.Close(op)
		return
	}
}

func init() {
	RootCmd.Flags().StringVarP(&flag.InUrl, "uri", "u", "", "Specify a target URI [e.g. -u https://example.com]")
	RootCmd.Flags().StringVarP(&flag.NetWork, "ip", "i", "", "Network segment [e.g. -i 192.168.1.0/24 or -i 192.168.1.1-10]")
	RootCmd.Flags().StringVarP(&flag.Mode, "mode", "m", "", "Specify the protocol [e.g. -m mysql/-m http]")
	RootCmd.Flags().StringVar(&flag.Type, "type", "", "Specify the type [e.g. --type tcp/--type udp]")
	RootCmd.Flags().IntVar(&flag.Timeout, "timeout", 5, "Response timeout time, the default is 5 seconds")
	RootCmd.Flags().IntVarP(&flag.Thread, "thread", "t", 500, "Number of concurrent threads")
	RootCmd.Flags().StringVarP(&flag.Port, "port", "p", "", "Custom scan ports [e.g. -p 80,443 or -p 1-65535]")
	RootCmd.Flags().StringVarP(&flag.Output, "output", "o", "output.txt", "Save the scan results to the specified file")
	RootCmd.Flags().StringVarP(&flag.File, "file", "f", "", "Parse the target from the specified file for batch recognition")
	RootCmd.Flags().BoolVar(&flag.NoIcmp, "np", false, "Not use ICMP/PING to detect surviving hosts")
	RootCmd.Flags().StringVarP(&flag.OutJson, "json", "j", "", "Scan result in json format [e.g. -j r.json]")
	RootCmd.Flags().BoolVar(&flag.NoColor, "nc", false, "Do not print character colors")
	RootCmd.Flags().IntVarP(&flag.Level, "level", "l", 3, "Specify log level (0:Fatal 1:Error 2:Info 3:Warning 4:Debug 5:Verbose)")
	RootCmd.Flags().StringVarP(&flag.Proxy, "proxy", "", "", "Use proxy scan, support http/socks5 protocol [e.g. --proxy socks5://127.0.0.1:1080]")
	RootCmd.Flags().BoolVarP(&flag.Pprof, "pprof", "d", false, "use pprof debug, on http://localhost:6060/debug/pprof/")
}

func Execute() {
	coloredcobra.Init(&coloredcobra.Config{
		RootCmd:         RootCmd,
		Headings:        coloredcobra.HiGreen + coloredcobra.Underline,
		Commands:        coloredcobra.Cyan + coloredcobra.Bold,
		Example:         coloredcobra.Italic,
		ExecName:        coloredcobra.Bold,
		Flags:           coloredcobra.Cyan + coloredcobra.Bold,
		NoExtraNewlines: true,
	})
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var RootCmd = &cobra.Command{
	Use: "dismap",
	Run: func(cmd *cobra.Command, args []string) {
		if flag.NetWork == "" && flag.File == "" && flag.InUrl == "" {
			configs.Banner()
			cmd.Help()
			return
		}
		_wg := &sync.WaitGroup{}
		if flag.Pprof {
			_wg.Add(1)
			go func() {
				logger.Info(http.ListenAndServe("localhost:6060", nil).Error())
			}()
		}

		configs.Banner()
		wg := &sync.WaitGroup{}
		lock := &sync.Mutex{}
		which(wg, lock)

		logger.Info("Identification completed and ended")

		_wg.Wait()
	},
}
