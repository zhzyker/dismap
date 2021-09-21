package dismap

import "flag"

type Options struct {
	Url     string
	Ips     string
	Ports   string
	File    string
	OutPut  string
	TimeOut int
	Threads int64
	NoIcmp  bool
}

func ParseOptions() *Options {
	p := &Options{}

	flag.StringVar(&p.Url, "url", "", "Specify a target URL [e.g. -url https://example.com]")
	flag.StringVar(&p.Ips, "ip", "", "Network segment [e.g. -ip 192.168.1.0/24 or -ip 192.168.1.1-10]")
	flag.StringVar(&p.Ports, "port", "", "Custom scan ports [e.g. -port 80,443 or -port 1-65535]")
	flag.StringVar(&p.File, "file", "", "Select a URL file for batch identification")
	flag.StringVar(&p.OutPut, "output", "output.txt", "Save the scan results to the specified file")
	flag.IntVar(&p.TimeOut, "timeout", 5, "Response timeout time, the default is 5 seconds")
	flag.Int64Var(&p.Threads, "thread", 256, "Number of concurrent threads, (adapted to two network segments 2x254)")
	flag.BoolVar(&p.NoIcmp, "np", false, "Not use ICMP/PING to detect surviving hosts")
	flag.Parse()

	return p
}
