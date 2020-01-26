package excavator

import (
	"github.com/godcong/excavator/net"
	"testing"
)

func TestExcavator_Run(t *testing.T) {
	excH := New(ActionArgs(RadicalTypeHanChengPinyin, RadicalTypeHanChengBihua, RadicalTypeHanChengBushou))
	e2 := excH.Run()
	if e2 != nil {
		t.Fatal(e2)
	}
	excK := New(ActionArgs(RadicalTypeKangXiPinyin, RadicalTypeKangXiBihua, RadicalTypeKangXiBushou))
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

func TestExcavator_SoList(t *testing.T) {
	debug = true
	soList := []string{"谫", "淯", "溧", "萋", "锯", "琦", "锬", "蕲", "誉", "蛛", "瘐", "洎", "锨", "莒", "誓", "脐", "淇", "琙", "膝", "胵", "遗", "逸", "镱", "隙", "薏", "蛭", "镆", "薤", "淹", "陷", "镥", "谳", "莳", "谏", "镜", "瘗", "萁", "淤", "险", "菔", "莅", "演", "禹", "蝜", "锂", "烯", "莺", "蝇", "荔", "燚", "营", "茯", "茵", "筵", "鄠", "萎", "藓", "谖", "隐", "理", "筮", "蜥", "熠", "蓟", "鄃", "莉", "萭", "锏", "淅", "浠", "棘", "蓰", "胰", "荫", "踽", "逶", "陴", "椅", "锱", "踞", "薇", "隈", "陎", "臆", "煜", "沿", "简", "葺", "渔", "茾", "踬", "银", "蝠", "铸", "荠", "蓣", "蜈", "楫", "邂", "腺", "瑜", "蛴", "鄅", "济", "腴", "漪", "胭", "谊", "茜", "镇", "葸", "琟", "谴", "谚", "萤", "隅", "熄", "锜", "蓄", "茱", "蹊", "萸", "谕", "茚", "院", "谐", "蜞", "蓖", "痿", "猬", "逝", "蝮", "脂", "鄢", "踯", "销", "漆", "蔚", "膣", "蜴", "焰", "猥", "蒍", "溪", "萦", "荐", "踺", "浴", "蒺", "锷", "谢", "篯", "琪", "谣", "谓", "滞", "渝", "蜘", "瘀", "葳", "邀", "镢", "滟", "遣", "胥", "遇", "瑀", "碛", "谥", "键", "蝓", "蓥", "铱", "蕙", "谦", "蒂", "镈", "燏", "镒", "鄞", "謇"}
	exc := New(ActionArgs(RadicalTypeHanChengSo, RadicalTypeKangXiSo))

	exc.SetSoList(soList)

	e := exc.Run()
	if e != nil {
		t.Fatal(e)
	}

}
