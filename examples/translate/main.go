// Package main demonstrates translating protovalidate constraint IDs to localized messages.
package main

import (
	"fmt"
	"log"

	"github.com/jzero-io/protovalidate-translator/translator"
)

func main() {
	// Example: translate a few constraint IDs (as returned by protovalidate violations)
	examples := []struct {
		ruleID string
		value  any
	}{
		{"float.lt", 100},
		{"string.min_len", 5},
		{"int32.gt", 0},
		{"string.email", nil},
	}

	fmt.Println("=== English (default) ===")
	for _, ex := range examples {
		data := map[string]any{}
		if ex.value != nil {
			data["Value"] = ex.value
		}
		msg, err := translator.TranslateDefault("en", ex.ruleID, data)
		if err != nil {
			log.Printf("translate %s: %v", ex.ruleID, err)
			continue
		}
		fmt.Printf("  %s -> %s\n", ex.ruleID, msg)
	}

	fmt.Println("\n=== 中文 ===")
	for _, ex := range examples {
		data := map[string]any{}
		if ex.value != nil {
			data["Value"] = ex.value
		}
		msg, err := translator.TranslateDefault("zh", ex.ruleID, data)
		if err != nil {
			log.Printf("translate %s: %v", ex.ruleID, err)
			continue
		}
		fmt.Printf("  %s -> %s\n", ex.ruleID, msg)
	}
}
