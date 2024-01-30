package internal

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zhzyker/dismap/configs"
	"github.com/zhzyker/dismap/pkg/logger"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "check rules",
	Run: func(cmd *cobra.Command, args []string) {
		for _, rule := range configs.RuleData {
			modes := strings.Split(rule.Mode, "|")
			types := strings.Split(rule.Type, "|")
			// Check the number of matches
			if rule.Mode == "" {
				if len(strings.Split(rule.Type, "|")) != 1 {
					logger.Error(fmt.Sprintf("Abnormal match pattern and quantity name: %-30v type: %-20v mode: %v", rule.Name, rule.Type, rule.Mode))
				}
			} else {
				if len(modes)+1 != len(types) {
					logger.Error(fmt.Sprintf("Abnormal match pattern and quantity name: %-30v type: %-20v mode: %v", rule.Name, rule.Type, rule.Mode))
				}

			}
			// check keyword
			for _, item := range types {
				if !(item == "body" || item == "header" || item == "ico") {
					logger.Error(fmt.Sprintf("Abnormal keyword, name: %-30v type: %-20v mode: %v", rule.Name, rule.Type, rule.Mode))
					break
				}
			}
			for _, item2 := range modes {
				if !(item2 == "" || item2 == "and" || item2 == "or") {
					logger.Error(fmt.Sprintf("Abnormal mode, name: %-30v type: %-20v mode: %v", rule.Name, rule.Type, rule.Mode))
					break
				}
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(checkCmd)
}
