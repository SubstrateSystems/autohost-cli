package utils

import "strings"

func ReplacePlaceholders(content string, values map[string]string) string {
	out := content
	for k, v := range values {
		out = strings.ReplaceAll(out, k, v)
	}
	return out
}
