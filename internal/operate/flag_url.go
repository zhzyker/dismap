package operate

import (
	"github.com/zhzyker/dismap/internal/output"
	"github.com/zhzyker/dismap/internal/parse"
	"github.com/zhzyker/dismap/internal/protocol"
	"github.com/zhzyker/dismap/pkg/logger"
	"os"
)

func FlagUrl(op *os.File, uri string, Args map[string]interface{}) {
	uri, scheme, host, port, err := parse.UriParse(uri)
	if logger.DebugError(err) {
		return
	}
	var res map[string]interface{}
	//Args["FlagMode"] = scheme
	switch scheme {
	case "http":
		res = protocol.DiscoverTcp(host, port, Args)
	case "https":
		res = protocol.DiscoverTls(host, port, Args)
	}
	//Args["FlagMode"] = ""
	parse.VerboseParse(res)
	output.Write(res, op)
}
