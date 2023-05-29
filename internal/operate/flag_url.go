package operate

import (
	"os"

	"github.com/zhzyker/dismap/internal/model"
	"github.com/zhzyker/dismap/internal/output"
	"github.com/zhzyker/dismap/internal/parse"
	"github.com/zhzyker/dismap/internal/protocol"
	"github.com/zhzyker/dismap/pkg/logger"
)

func FlagUrl(op *os.File, uri string) {
	uri, scheme, host, port, err := parse.UriParse(uri)
	if logger.DebugError(err) {
		return
	}
	var res *model.Result
	//Args["FlagMode"] = scheme
	switch scheme {
	case "http":
		res = protocol.DiscoverTcp(host, port)
	case "https":
		res = protocol.DiscoverTls(host, port)
	}
	//Args["FlagMode"] = ""
	parse.VerboseParse(res)
	output.Write(res, op)
}
