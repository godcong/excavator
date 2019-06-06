package excavator

import "testing"

// TestCommonlyBase ...
func TestCommonlyBase(t *testing.T) {
	Debug()
	rc := RootRadicalCharacter{Class: 0, Character: "勺", Link: "/hans/勺", Pinyin: []string{"sháo"}}
	CommonlyBase("http://www.zdic.net", &rc)
}
