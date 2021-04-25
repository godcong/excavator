package excavator

import (
	"excavator/config"
	"excavator/models"
	"excavator/net"
	"fmt"
	"testing"
	"time"
)

func TestExcavator_Run(t *testing.T) {
	before := time.Now()

	excH := New()

	excH.action = config.ActionParse

	e2 := excH.Run()
	if e2 != nil {
		t.Fatal(e2)
	}

	fmt.Printf("took %v\n", time.Since(before))
}

func TestGetCharacter(t *testing.T) {

	before := time.Now()

	excH := New()

	//t.Log(e)
	debug = true
	soList := []string{}
	// soList = append(soList, *regular.GetRegularList()...)
	// soList = append(soList, []string{"溧", "萋", "锯", "琦", "锬", "蕲", "誉", "蛛", "瘐", "洎", "莒", "誓", "脐", "淇", "膝", "遗", "逸", "隙", "薏", "蛭", "镆", "薤"}...)
	// soList = append(soList, []string{"淹", "陷", "谳", "莳", "谏", "镜", "瘗", "萁", "淤", "险", "菔", "莅", "演", "禹", "烯", "莺", "蝇", "荔", "营", "茯", "茵", "筵"}...)
	// soList = append(soList, []string{"萎", "藓", "谖", "隐", "理", "筮", "蜥", "熠", "蓟", "莉", "淅", "浠", "棘", "蓰", "胰", "荫", "踽", "逶", "陴", "椅", "锱", "踞"}...)
	// soList = append(soList, []string{"薇", "隈", "陎", "臆", "煜", "沿", "葺", "渔", "茾", "踬", "银", "蝠", "铸", "荠", "蓣", "蜈", "楫", "邂", "瑜", "蛴", "济", "腴"}...)
	// soList = append(soList, []string{"漪", "胭", "谊", "茜", "镇", "葸", "谴", "谚", "萤", "隅", "熄", "蓄", "茱", "蹊", "萸", "谕", "茚", "院", "谐", "蜞", "痿", "猬"}...)
	// soList = append(soList, []string{"逝", "蝮", "脂", "鄢", "踯", "销", "漆", "蔚", "膣", "蜴", "焰", "猥", "溪", "萦", "荐", "踺", "浴", "蒺", "锷", "谢", "琪", "谣"}...)
	// soList = append(soList, []string{"谓", "滞", "渝", "蜘", "瘀", "葳", "邀", "镢", "遣", "胥", "遇", "碛", "谥", "键", "蝓", "蓥", "蕙", "谦", "蒂", "镒", "鄞", "謇"}...)
	// soList = append(soList, []string{"灶", "朎", "犯", "蓖", "火", "迂", "镈", "邚", "趺", "讥", "込", "玖", "汊", "言", "灵", "矽", "疒", "镢", "趸", "艺", "钎", "阝"}...)
	// soList = append(soList, []string{"有", "虫", "邜", "辻", "忛", "竺", "玑", "石", "虵", "锜", "肭", "钍", "䇖", "简", "腺", "迅", "狄", "锂", "足", "王", "镥", "钓"}...)
	// soList = append(soList, []string{"苄", "犹", "笃", "铱", "疖", "月", "术", "锨", "趷", "记", "辽", "迁", "疚", "篯", "肌", "钆", "迄", "镱", "疗", "朱", "虬", "虲"}...)
	// soList = append(soList, []string{"疜", "邝", "趼", "䟞", "辺", "朋", "计", "邓", "犸", "犷", "边", "艹", "訊", "让", "汇", "矶", "朰", "竹", "蒍", "锏", "笂", "节"}...)
	// soList = append(soList, []string{"谫", "钇", "䖝", "讫", "䟖", "认", "笁", "矴", "滟", "㬳", "钗", "朑", "㐀", "岩", "怅", "𪇭", "黄", "和", "木", "华", "鼹", "齑"}...)
	// soList = append(soList, []string{"悵", "遥", "朐", "王", "音", "森", "雪", "盐", "张", "阧", "障", "陇", "階", "阾", "降", "隄", "陞", "陈", "隯", "隗", "阱", "陲"}...)
	// soList = append(soList, []string{"阼", "阬", "陳", "陌", "隀", "阠", "陿", "隉", "陉", "陧", "陱", "阶", "陂", "陽", "陪", "陗", "阪", "阞", "隭", "阬", "陆", "隡"}...)
	// soList = append(soList, []string{"陦", "阳", "阽", "陔", "隄", "防", "隕", "陨", "隬", "隣", "陕", "隝", "隢", "陹", "限", "附", "陰", "阮", "隚", "陜", "陸", "阦"}...)
	// soList = append(soList, []string{"阞", "隘", "隴", "隍", "陶", "隀", "隔", "陀", "陬", "隨", "除", "阩", "阴", "陋", "队", "隳", "隊", "隋", "陘", "阥", "隆", "陡"}...)
	// soList = append(soList, []string{"隧", "陯", "隍", "陛", "阻", "隯", "陃", "隘", "阨", "阨", "隌", "陔", "陵", "附", "陟", "队", "随", "隌", "陃", "喦", "黽", "乾"}...)
	soList = append(soList, []string{}...)

	invalidList := []string{}
	invalidList = append(invalidList, []string{"疓", "忔", "朏", "肊", "朒", "㶣", "禸", "䟕", "趽", "汈", "朲", "䟔", "忕", "虳", "矷", "㶡", "迆", "㽱", "犻", "氻"}...)
	invalidList = append(invalidList, []string{"訉", "䟛", "䚰", "玌", "虴", "疛", "邒", "灴", "趹", "朳", "阞", "邘", "㽲", "瑀", "燏", "琟", "鄅", "鄃", "萭", "蝜"}...)
	invalidList = append(invalidList, []string{"燚", "鄠", "淯", "琙", "胵", "陾", "隓", "隂", "陖", "陠", "隩", "隥", "陁", "陼", "陊", "陑", "陫", "陚", "阺", "阫"}...)
	invalidList = append(invalidList, []string{"阭", "陏", "阷", "隫", "隞", "陒", "阤", "隤", "陮", "陙", "阸", "隑", "陖"}...)
	invalidList = append(invalidList, []string{"氧", "挡", "眉"}...)
	invalidList = append(invalidList, []string{}...)

	invalidMap := map[string]bool{}

	for _, invalid := range invalidList {
		invalidMap[invalid] = true
	}

	soList = append(soList, []string{"签"}...)

	invalids := []string{}

	for _, so := range soList {
		if invalidMap[so] {
			continue
		}

		rc := models.HanChengChar{
			Unid: int(([]rune(so))[0]),
		}
		has, err := excH.Db.Get(&rc)

		if err != nil {
			panic(err)
		}

		if !has {
			fmt.Printf("%s url没有记录\n", so)
			invalids = append(invalids, so)
			continue
		}

		html_node, e := net.CacheQuery(UrlMerge(excH.base_url, rc.Url))

		if e != nil {
			panic(e)
		}

		e = getCharacter(excH, rc.Unid, html_node)

		if e != nil {
			panic(e)
		}
	}

	fmt.Println(invalids)

	fmt.Println(len(soList))

	fmt.Println(len(invalidList))

	fmt.Printf("took %v\n", time.Since(before))
}
