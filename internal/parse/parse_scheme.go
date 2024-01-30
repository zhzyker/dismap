package parse

import (
	"fmt"
	"strconv"

	"github.com/zhzyker/dismap/internal/model"
)

func SchemeParse(result *model.Result) string {
	path := result.Path
	scheme := result.Protocol
	port := result.Port
	host := result.Host
	if scheme != "" && path != "" {
		result.Uri = fmt.Sprintf("%s://%s:%s%s", scheme, host, strconv.Itoa(port), path)
		return result.Uri
	} else if scheme != "" {
		result.Uri = fmt.Sprintf("%s://%s:%s", scheme, host, strconv.Itoa(port))
		return result.Uri
	} else {
		result.Uri = fmt.Sprintf("%s:%s", host, strconv.Itoa(port))
		return result.Uri
	}
}
