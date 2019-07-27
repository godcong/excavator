package excavator

import (
	"regexp"
	"strings"
)

// StringClearUp ...
func StringClearUp(s string) (t string) {
	t = strings.ReplaceAll(s, "\n", "")
	t = strings.TrimSpace(t)
	rgx := regexp.MustCompile(`[\s]+`)
	return rgx.ReplaceAllString(t, " ")
}
