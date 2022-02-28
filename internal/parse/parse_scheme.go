package parse

import (
	"fmt"
	"strconv"
)

func SchemeParse(result map[string]interface{}) string {
	path := result["path"].(string)
	scheme := result["protocol"].(string)
	port := result["port"].(int)
	host := result["host"].(string)
	if scheme != "" && path != "" {
		result["uri"] = fmt.Sprintf("%s://%s:%s%s",scheme, host, strconv.Itoa(port), path)
		return result["uri"].(string)
	} else if scheme != "" {
		result["uri"] = fmt.Sprintf("%s://%s:%s",scheme, host, strconv.Itoa(port))
		return result["uri"].(string)
	} else {
		result["uri"] = fmt.Sprintf("%s:%s", host, strconv.Itoa(port))
		return result["uri"].(string)
	}
}
