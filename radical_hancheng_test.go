package excavator

import "testing"

func TestRadical_Hancheng(t *testing.T) {
	db = InitMysql("localhost:3306", "root", "111111")
	db.ShowSQL()
	e := db.Sync2(&RadicalCharacter{})
	if e != nil {
		panic(e)
	}
	grabRadicalList(SearchTypeBushou, "http://hy.httpcn.com/bushou/zi/")

	//grabRadicalList(SearchTypePinyin, "http://hy.httpcn.com/pinyin/kangxi/")
}
