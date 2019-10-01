package excavator

import "testing"

func TestRadical_Hancheng(t *testing.T) {
	db = InitMysql("localhost:3306", "root", "111111")
	e := db.Sync2(&RadicalCharacter{})
	if e != nil {
		panic(e)
	}
	grabRadicalList("http://hy.httpcn.com/bushou/zi/")
}
