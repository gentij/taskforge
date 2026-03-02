package utils

import (
	"strings"

	"github.com/muesli/reflow/wordwrap"
)

func WrapText(text string, width int) string {
	if width <= 0 {
		return text
	}
	lines := strings.Split(text, "\n")
	wrapped := make([]string, 0, len(lines))
	for _, line := range lines {
		wrapped = append(wrapped, wordwrap.String(line, width))
	}
	return strings.Join(wrapped, "\n")
}
