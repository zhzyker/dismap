package rule

import (
	"encoding/json"
	"log"
	"os"
	"testing"
)

func TestDumpRule(t *testing.T) {
	data, err := json.Marshal(RuleData)
	if err != nil {
		log.Fatal(err)
	}
	f, err := os.Create("rule.json")
	if err != nil {
		log.Fatal(err)
	}
	f.Write(data)
}
