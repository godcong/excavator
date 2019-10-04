package excavator

import "testing"

func TestExcavator_Run(t *testing.T) {
	exc := New(RadicalTypeKangXiPinyin, SkipArgs("radical"))
	e := exc.Run()
	if e != nil {
		return
	}
}
