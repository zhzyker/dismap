package match

import (
	"encoding/json"
	"github.com/zhzyker/dismap/pkg/logger"
	"github.com/zhzyker/dismap/pkg/requet/http"
)

// addMatch 结构体定义了识别结果的标识符
type addMatch struct {
	I string `json:"identify"`
}

// identifyResult 用于生成 JSON 格式的指纹识别结果
func identifyResult(matches []addMatch, res http.Responses) []byte {
	// 将结果转化为 JSON 格式
	jsonResult, err := json.Marshal(map[string]interface{}{
		"url":     res.Url,
		"matches": matches,
		"title":   res.Title,
	})
	if err != nil {
		logger.ERR("Failed to marshal JSON result: " + err.Error())
		return nil
	}
	return jsonResult
}
