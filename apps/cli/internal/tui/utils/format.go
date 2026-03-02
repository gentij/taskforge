package utils

import (
	"encoding/json"
	"strings"
)

func Truncate(value string, width int) string {
	if width <= 0 {
		return ""
	}
	runes := []rune(value)
	if len(runes) <= width {
		return value
	}
	if width <= 3 {
		return string(runes[:width])
	}
	return string(runes[:width-3]) + "..."
}

func Indent(text string, prefix string) string {
	lines := strings.Split(text, "\n")
	for i := range lines {
		lines[i] = prefix + lines[i]
	}
	return strings.Join(lines, "\n")
}

func PrettyJSON(raw string) string {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return "{}"
	}
	var payload any
	if err := json.Unmarshal([]byte(trimmed), &payload); err != nil {
		return raw
	}
	formatted, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return raw
	}
	return string(formatted)
}

func FilterLines(content string, query string) string {
	query = strings.TrimSpace(strings.ToLower(query))
	if query == "" {
		return content
	}
	lines := strings.Split(content, "\n")
	filtered := make([]string, 0, len(lines))
	for _, line := range lines {
		if strings.Contains(strings.ToLower(line), query) {
			filtered = append(filtered, line)
		}
	}
	if len(filtered) == 0 {
		return "No matches"
	}
	return strings.Join(filtered, "\n")
}
