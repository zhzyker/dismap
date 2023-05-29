package parse

import (
	"encoding/hex"
	"fmt"

	"github.com/zhzyker/dismap/internal/model"
	"github.com/zhzyker/dismap/pkg/logger"
)

func VerboseParse(res *model.Result) {
	logger.Verbose(fmt.Sprintf("Hex dump\n%s", hex.Dump(res.BannerB)))
	r := "\n"
	// for k, v := range res {
	// 	r += fmt.Sprintf("%18s: %s", fmt.Sprintf(k), fmt.Sprintln(v))
	// }
	logger.Verbose(fmt.Sprintf("Dismap identify result\n%s", r))
}
