package excavator

import (
	"strings"
)

func StringSplite(base, sta, end string) []string {
	sa := make([]string, 0)
	size := len(base)
	for i := 0; i < size; {
		s := strings.Index(base, sta)
		if s < 0 {
			break
		}
		base = base[s:]
		e := strings.Index(base, end)
		if e < 0 {
			break
		}
		e = e + len(end)
		sa = append(sa, base[0:e])

		i += e
		base = base[e:]
	}
	return sa
}
