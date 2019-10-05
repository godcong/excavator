package excavator

import "testing"

func TestExcavator_Run(t *testing.T) {
	excK := New(RadicalTypeKangXiPinyin, ActionArgs(RadicalTypeKangXiPinyin, RadicalTypeKangXiBihua, RadicalTypeKangXiBushou))
	e1 := excK.Run()
	if e1 != nil {
		t.Fatal(e1)
	}
	excH := New(RadicalTypeHanChengPinyin, ActionArgs(RadicalTypeHanChengPinyin, RadicalTypeHanChengBihua, RadicalTypeHanChengBushou))
	e2 := excH.Run()
	if e2 != nil {
		t.Fatal(e2)
	}
}
