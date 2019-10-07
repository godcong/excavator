# excavator


How to use
```
	excH := New(RadicalTypeHanChengPinyin, ActionArgs(RadicalTypeHanChengPinyin, RadicalTypeHanChengBihua, RadicalTypeHanChengBushou))
	e2 := excH.Run()
	if e2 != nil {
		t.Fatal(e2)
	}
	excK := New(RadicalTypeKangXiPinyin, ActionArgs(RadicalTypeKangXiPinyin, RadicalTypeKangXiBihua, RadicalTypeKangXiBushou))
	e1 := excK.Run()
	if e1 != nil {
		t.Fatal(e1)
	}
```
if want change db address use:
```
   exc :=New(RadicalTypeHanChengPinyin, DBArgs(#db#))
   or
   exc.SetDB(#db#)
```