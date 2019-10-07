package excavator

import (
	"github.com/godcong/excavator/net"
	"testing"
)

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

func TestGetCharacter(t *testing.T) {
	debug = true
	document, e := net.CacheQuery("http://hy.httpcn.com/html/kangxi/35/KOILPWUYXVXVUYDPW")
	if e != nil {
		t.Fatal(e)
	}
	character := getCharacter(document, &RadicalCharacter{}, true)
	log.Info(character)
	//http://hy.httpcn.com/html/kangxi/35/KOILPWUYXVXVUYDPW
}
