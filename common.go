package excavator

import (
	"regexp"
	"strings"
)

// StringClearUp ...
func StringClearUp(s string) (t string) {
	//t = strings.ReplaceAll(s, "\n", "")
	t = strings.TrimSpace(t)
	rgx := regexp.MustCompile(`[\s]+`)
	return rgx.ReplaceAllString(t, " ")
}

/*URL 拼接地址 */
func URL(prefix string, uris ...string) string {
	end := len(prefix)
	if end > 1 && prefix[end-1] == '/' {
		prefix = prefix[:end-1]
	}

	var url = []string{prefix}
	for _, v := range uris {
		url = append(url, TrimSlash(v))
	}
	return strings.Join(url, "/")
}

// TrimSlash ...
func TrimSlash(s string) string {
	if size := len(s); size > 1 {
		if s[size-1] == '/' {
			s = s[:size-1]
		}
		if s[0] == '/' {
			s = s[1:]
		}
	}
	return s
}

func ContainsQuery(input []string, query string) bool {
	for _, str := range input {
		if strings.Contains(str, query) && strings.Contains(str, "）") && strings.Contains(str, "（") {
			return true
		}
	}

	return false
}
