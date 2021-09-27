package lib

import "flag"

var NetWork string
var InUrl string
var TimeOut int
var Thread int
var Port string
var OutPut string
var Files string
var NoIcmp bool

func init() {
	flag.StringVar(&InUrl, "url", "", "Specify a target URL [e.g. -url https://example.com]")
	flag.StringVar(&NetWork,"ip", "", "Network segment [e.g. -ip 192.168.1.0/24 or -ip 192.168.1.1-10]")
	flag.IntVar(&TimeOut, "timeout", 5, "Response timeout time, the default is 5 seconds")
	flag.IntVar(&Thread, "thread", 508, "Number of concurrent threads, (adapted to two network segments 2x254)")
	flag.StringVar(&Port, "port", "", "Custom scan ports [e.g. -port 80,443 or -port 1-65535]")
	flag.StringVar(&OutPut, "output", "output.txt", "Save the scan results to the specified file")
	flag.StringVar(&Files, "file", "", "Select a URL file for batch identification")
	flag.BoolVar(&NoIcmp, "np", false, "Not use ICMP/PING to detect surviving hosts")
}

