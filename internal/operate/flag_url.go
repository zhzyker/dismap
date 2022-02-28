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
	if Args["FlagMode"] == Args["FlagType"] {
		if scheme == "https" {
			Args["FlagType"] = "tls"
			Args["FlagMode"] = scheme
		} else {
			Args["FlagType"] = "tcp"
			Args["FlagMode"] = scheme
		}
	}
	res := protocol.Discover(host, port, Args)
	if Args["FlagMode"] == Args["FlagType"] {
		Args["FlagType"] = ""
		Args["FlagMode"] = ""
	}
	parse.VerboseParse(res)
	output.Write(res, op)
}
