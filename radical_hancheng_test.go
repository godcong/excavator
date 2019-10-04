package excavator

import "testing"

func TestRadical_Hancheng(t *testing.T) {
	db = InitMysql("localhost:3306", "root", "111111")
	db.ShowSQL()
	e := db.Sync2(&RadicalCharacter{})
	if e != nil {
		panic(e)
	}
	exc := New(RadicalTypeKangXiPinyin)
	exc.Run()
	//grabRadicalList(RadicalTypeHanChengBushou, DefaultMainPage)

	//grabRadicalList(RadicalTypeHanChengPinyin, DefaultMainPage)

	//grabRadicalList(RadicalTypeHanChengBihua, DefaultMainPage)

	//grabRadicalList(RadicalTypeKangXiBushou, DefaultMainPage)

	//grabRadicalList(RadicalTypeKangXiPinyin, DefaultMainPage)

	//grabRadicalList(RadicalTypeKangXiBihua, DefaultMainPage)
}
