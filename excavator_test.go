package excavator

import "testing"

func TestExcavator_Run(t *testing.T) {
	exc := New(RadicalTypeKangXiPinyin, ActionArgs(RadicalTypeKangXiPinyin, RadicalTypeKangXiBihua, RadicalTypeKangXiBushou))
	e := exc.Run()
	if e != nil {
		return
	}

}
