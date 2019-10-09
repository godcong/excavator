# excavator

## How to use
```go
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
	
	this tool used mysql address localhost:3306/root/111111 with default database.	
	
	if want change db address use:
	
```go
   exc :=New(RadicalTypeHanChengPinyin, DBArgs(#db#))
   //or
   exc.SetDB(#db#)
```
if you want to change the database mssql/postgre...
new a engine with xorm by yourself 
then add it by DBArgs().
