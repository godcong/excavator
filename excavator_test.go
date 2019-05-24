package excavator_test

import (
	"strings"
	"testing"
)

func TestRadical_Add(t *testing.T) {
	text := "汉字五行：土　是否为常用字：否"
	s := strings.SplitAfter(text, "：")
	t.Log(s)
}
