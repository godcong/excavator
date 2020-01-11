package excavator

import (
	"github.com/godcong/excavator/net"
	"testing"
)

func TestExcavator_Run(t *testing.T) {
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
}

func TestGetCharacter(t *testing.T) {
	db := InitMysql("localhost:3306", "root", "111111")
	e := db.Sync(&Character{})
	t.Log(e)
	debug = true
	document, e := net.CacheQuery("http://hy.httpcn.com/html/kangxi/26/PWUYUYUYPWTBKODAZ/")
	if e != nil {
		t.Fatal(e)
	}
	character := getCharacter(document, &RadicalCharacter{Zi: "㰄"}, true)
	log.Info(character)
	i, e := character.InsertOrUpdate(db.Where(""))
	t.Log(i, e)
}
