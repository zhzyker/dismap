package flag

var (
	NetWork string
	InUrl   string
	Timeout int
	Thread  int
	Port    string
	Output  string
	File    string
	NoIcmp  bool
	NoColor bool
	Mode    string
	Type    string
	Help    bool
	Level   int
	Proxy   string
	OutJson string
	Pprof   bool
)

// func init() {
// 	flag.StringVarP(&InUrl, "uri", "u", "", "Specify a target URI [e.g. -u https://example.com]")
// 	flag.StringVarP(&NetWork,"ip", "i", "", "Network segment [e.g. -i 192.168.1.0/24 or -i 192.168.1.1-10]")
// 	flag.StringVarP(&Mode,"mode", "m", "", "Specify the protocol [e.g. -m mysql/-m http]")
// 	flag.StringVar(&Type,"type", "", "Specify the type [e.g. --type tcp/--type udp]")
// 	flag.IntVar(&Timeout, "timeout", 5, "Response timeout time, the default is 5 seconds")
// 	flag.IntVarP(&Thread, "thread", "t", 500, "Number of concurrent threads")
// 	flag.StringVarP(&Port, "port", "p", "", "Custom scan ports [e.g. -p 80,443 or -p 1-65535]")
// 	flag.StringVarP(&Output, "output", "o", "output.txt", "Save the scan results to the specified file")
// 	flag.StringVarP(&File, "file", "f", "", "Parse the target from the specified file for batch recognition")
// 	flag.BoolVar(&NoIcmp, "np",false, "Not use ICMP/PING to detect surviving hosts")
// 	flag.StringVarP(&OutJson, "json", "j", "", "Scan result in json format [e.g. -j r.json]")
// 	flag.BoolVar(&NoColor, "nc", false, "Do not print character colors")
// 	flag.IntVarP(&Level, "level", "l", 3, "Specify log level (0:Fatal 1:Error 2:Info 3:Warning 4:Debug 5:Verbose)")
// 	flag.StringVarP(&Proxy, "proxy", "", "", "Use proxy scan, support http/socks5 protocol [e.g. --proxy socks5://127.0.0.1:1080]")
// 	flag.BoolVarP(&Help, "help", "h",false, "Show help")
// }

// func Flags() map[string]interface{} {
// 	flag.Parse()
// 	if Help {
// 		flag.PrintDefaults()
// 		os.Exit(0)
// 	}
// 	flags := map[string]interface{}{
// 		"FlagUrl":       InUrl,
// 		"FlagNetwork":   NetWork,
// 		"FlagMode":      Mode,
// 		"FlagType":      Type,
// 		"FlagTimeout":   Timeout,
// 		"FlagThread":    Thread,
// 		"FlagPort":      Port,
// 		"FlagOutput":    Output,
// 		"FlagFile":      File,
// 		"FlagNoIcmp":    NoIcmp,
// 		"FlagProxy":	 Proxy,
// 		"FlagOutJson":	 OutJson,
// 	}
// 	return flags
// }
