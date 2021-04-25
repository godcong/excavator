package excavator

import (
	"errors"
	"excavator/models"
	"math/bits"
	"strconv"
	"strings"

	"github.com/antchfx/htmlquery"
	"golang.org/x/exp/utf8string"
	"golang.org/x/net/html"
)

//康熙字典 & 姓名学笔画
func parseKangXi(exc *Excavator, unid int, glyph *models.Glyph, html_node *html.Node, stroke_kang map[rune]int) (err error) {
	science_stroke := models.ScienceStroke{
		Unid: unid,
	}

	if len(stroke_kang) > 0 {
		if kang_stroke, has := stroke_kang[rune(unid)]; has {
			if len(stroke_kang) > 1 {
				panic("输入的康熙笔画异常")
			}

			science_stroke.ScienceStroke = kang_stroke

			parseDictComment(exc, unid, html_node, false)
		} else {
			if len(stroke_kang) > 1 {
				kang_stroke, _ := parseDictComment(exc, unid, html_node, true)
				science_stroke.ScienceStroke = kang_stroke
			} else {
				for _, my_stroke := range stroke_kang {
					science_stroke.ScienceStroke = my_stroke
				}
			}
		}

	} else {
		//当康熙笔画列表中为空
		if glyph.TraditionalTotalStroke == 0 {
			if glyph.Stroke == 0 {
				panic("非康熙字部首笔画异常")
			}
			science_stroke.ScienceStroke = glyph.Stroke
		} else {
			science_stroke.ScienceStroke = glyph.TraditionalTotalStroke
		}
	}

	err = InsertOrUpdate(exc.Db, &science_stroke)
	if err != nil {
		panic(err)
	}

	return nil
}

func parseZiBuShou(glyph *models.Glyph, he_xin_block *html.Node, ji_ben_block *html.Node) (err error) {
	radical_stroke := htmlquery.FindOne(he_xin_block, ".//span[contains(@class, 'b') and contains(text(), '部首：') and not(contains(text(), '简')) and not(contains(text(), '繁'))]/following-sibling::span[1]/following-sibling::text()[1]")

	if radical_stroke != nil {
		radical_stroke_str := strings.TrimSpace(htmlquery.InnerText(radical_stroke))
		radical_stroke_num, err := strconv.ParseUint(radical_stroke_str, 10, bits.UintSize)
		if err != nil {
			panic(err)
		}
		glyph.RadicalStroke = int(radical_stroke_num)
	}

	stroke := htmlquery.FindOne(he_xin_block, ".//span[contains(@class, 'b') and contains(text(), '部首：') and not(contains(text(), '简')) and not(contains(text(), '繁'))]/following-sibling::span[2]/following-sibling::text()[1]")

	if stroke != nil {
		stroke_num_str := strings.TrimSpace(htmlquery.InnerText(stroke))
		stroke_num, err := strconv.ParseUint(stroke_num_str, 10, bits.UintSize)
		if err != nil {
			panic(err)
		}
		glyph.Stroke = int(stroke_num)
	} else {
		ji_ben_stroke := htmlquery.FindOne(ji_ben_block, "./text()[contains(., '笔画数：')]")

		if ji_ben_stroke != nil {
			ji_ben_stroke_str := strings.TrimSpace(htmlquery.InnerText(ji_ben_stroke))
			ji_ben_stroke_str_num := strings.TrimSpace(strings.Split(strings.Split(ji_ben_stroke_str, "笔画数：")[1], "；")[0])

			ji_ben_stroke_num, err := strconv.ParseUint(ji_ben_stroke_str_num, 10, bits.UintSize)
			if err != nil {
				panic(err)
			}
			glyph.Stroke = int(ji_ben_stroke_num)
		}

	}

	radical := htmlquery.FindOne(he_xin_block, ".//span[contains(@class, 'b') and contains(text(), '部首：') and not(contains(text(), '简')) and not(contains(text(), '繁'))]/following-sibling::text()[1]")

	if radical != nil {
		glyph.Radical = strings.TrimSpace(htmlquery.InnerText(radical))
		if len(glyph.Radical) == 0 {
			if len(glyph.ShouWei) == 0 {
				return errors.New("没有部首信息")
			}
			glyph.Radical = string(([]rune(glyph.ShouWei))[0])
		}
	} else {
		ji_ben_radical := htmlquery.FindOne(ji_ben_block, "./text()[contains(., '部首：') and not(contains(., '难检字'))]")

		if ji_ben_radical != nil {
			ji_ben_radical_str := strings.TrimSpace(htmlquery.InnerText(ji_ben_radical))
			ji_ben_radical_str_inner := strings.TrimSpace(strings.Split(strings.Split(ji_ben_radical_str, "部首：")[1], "；")[0])

			glyph.Radical = ji_ben_radical_str_inner
		}
	}

	return nil
}

func parseSimple(glyph *models.Glyph, he_xin_block *html.Node) {
	simplified_radical := htmlquery.FindOne(he_xin_block, ".//span[contains(@class, 'b') and contains(text(), '简体部首：')]/following-sibling::text()[1]")

	if simplified_radical != nil {
		glyph.SimplifiedRadical = strings.TrimSpace(htmlquery.InnerText(simplified_radical))
		if !strings.Contains(glyph.SimplifiedRadical, "难检字") {
			glyph.SimplifiedRadical = ""
		}
	}

	simplified_radical_stroke := htmlquery.FindOne(he_xin_block, ".//span[contains(@class, 'b') and contains(text(), '简体部首：')]/following-sibling::span[1]/following-sibling::text()[1]")

	if simplified_radical_stroke != nil {
		simplified_radical_stroke_num_str := strings.TrimSpace(htmlquery.InnerText(simplified_radical_stroke))
		if len(simplified_radical_stroke_num_str) > 0 {
			simplified_radical_stroke_num, err := strconv.ParseUint(simplified_radical_stroke_num_str, 10, bits.UintSize)
			if err != nil {
				panic(err)
			}
			glyph.SimplifiedRadicalStroke = int(simplified_radical_stroke_num)
		}
	}

	simplified_total_stroke := htmlquery.FindOne(he_xin_block, ".//span[contains(@class, 'b') and contains(text(), '简体部首：')]/following-sibling::span[2]/following-sibling::text()[1]")

	if simplified_total_stroke != nil {
		simplified_total_stroke_num, err := strconv.ParseUint(strings.TrimSpace(htmlquery.InnerText(simplified_total_stroke)), 10, bits.UintSize)
		if err != nil {
			panic(err)
		}
		glyph.SimplifiedTotalStroke = int(simplified_total_stroke_num)
	}
}

func parseTraditional(glyph *models.Glyph, he_xin_block *html.Node) {
	traditional_radical := htmlquery.FindOne(he_xin_block, ".//span[contains(@class, 'b') and contains(text(), '繁体部首：')]/following-sibling::text()[1]")

	if traditional_radical != nil {
		glyph.TraditionalRadical = strings.TrimSpace(htmlquery.InnerText(traditional_radical))
	}

	traditional_radical_stroke := htmlquery.FindOne(he_xin_block, ".//span[contains(@class, 'b') and contains(text(), '繁体部首：')]/following-sibling::span[1]/following-sibling::text()[1]")

	if traditional_radical_stroke != nil {
		traditional_radical_stroke_num, err := strconv.ParseUint(strings.TrimSpace(htmlquery.InnerText(traditional_radical_stroke)), 10, bits.UintSize)
		if err != nil {
			panic(err)
		}
		glyph.TraditionalRadicalStroke = int(traditional_radical_stroke_num)
	}

	traditional_total_stroke := htmlquery.FindOne(he_xin_block, ".//span[contains(@class, 'b') and contains(text(), '繁体部首：')]/following-sibling::span[2]/following-sibling::text()[1]")

	if traditional_total_stroke != nil {
		traditional_total_stroke_num, err := strconv.ParseUint(strings.TrimSpace(htmlquery.InnerText(traditional_total_stroke)), 10, bits.UintSize)
		if err != nil {
			panic(err)
		}
		glyph.TraditionalTotalStroke = int(traditional_total_stroke_num)
	}
}

//字形
func parseZiCharacter(exc *Excavator, glyph *models.Glyph, unid int, html_node *html.Node, he_xin_block *html.Node, ji_ben_block *html.Node) (err error) {
	zi_xing_block := htmlquery.FindOne(html_node, "//div[contains(@class, 'text16')]/span[contains(text(), '字形')]/..")

	if zi_xing_block == nil {
		panic("字形结构区找不到")
	}

	shou_wei := htmlquery.FindOne(zi_xing_block, ".//span[contains(text(), '分')]/following-sibling::text()[1]")

	if shou_wei != nil {
		shou_wei_str := strings.TrimSpace(htmlquery.InnerText(shou_wei))
		glyph.ShouWei = strings.TrimSpace(strings.Split(strings.Split(shou_wei_str, "]：")[1], "[")[0])
	}

	bu_jian := htmlquery.FindOne(zi_xing_block, ".//span[contains(text(), '件')]/following-sibling::text()[1]")

	if bu_jian != nil {
		bu_jian_str := strings.TrimSpace(htmlquery.InnerText(bu_jian))
		glyph.BuJian = strings.TrimSpace(strings.Split(bu_jian_str, "]：")[1])
	}

	bi_shun := htmlquery.FindOne(zi_xing_block, ".//span[contains(text(), '号')]/following-sibling::text()[1]")

	if bi_shun != nil {
		bi_shun_str := strings.TrimSpace(htmlquery.InnerText(bi_shun))
		glyph.BiHao = strings.TrimSpace(strings.Split(bi_shun_str, "]：")[1])
	}

	bi_du := htmlquery.FindOne(zi_xing_block, ".//span[contains(text(), '写')]/following-sibling::text()[1]")

	if bi_du != nil {
		bi_du_str := strings.TrimSpace(htmlquery.InnerText(bi_du))
		glyph.BiDu = strings.TrimSpace(strings.Split(bi_du_str, "]：")[1])
	}

	parseSimple(glyph, he_xin_block)

	parseTraditional(glyph, he_xin_block)

	err = parseZiBuShou(glyph, he_xin_block, ji_ben_block)

	ji_ben_radical_is := htmlquery.FindOne(ji_ben_block, "./text()[contains(., '作偏旁') or contains(., '汉字部首')]")

	glyph.AsRadical = false

	if ji_ben_radical_is != nil {
		glyph.AsRadical = true
	} else {
		unid_char := string(rune(unid))
		if len(glyph.Radical) > 0 {
			if unid_char == glyph.Radical {
				glyph.AsRadical = true
			}
		} else {
			if glyph.SimplifiedTotalStroke > 0 || glyph.TraditionalTotalStroke > 0 || glyph.Stroke > 0 {
				if unid_char == glyph.SimplifiedRadical || unid_char == glyph.TraditionalRadical || unid_char == glyph.Radical {
					glyph.AsRadical = true
				}
			} else {
				panic("没有部首信息")
			}
		}
	}

	err = InsertOrUpdate(exc.Db, glyph)
	if err != nil {
		panic(err)
	}

	return nil
}

func parseRegular(wu_xing_ji_xiong_regular_str string) (chang_yong string) {
	if strings.Contains(wu_xing_ji_xiong_regular_str, "常用字：") {
		chang_yong = utf8string.NewString(strings.TrimSpace(strings.Split(wu_xing_ji_xiong_regular_str, "常用字：")[1])).Slice(0, 1)
	} else {
		chang_yong = ""
	}

	return chang_yong
}

func parseNameScience(min_su *models.MinSu, min_su_block *html.Node) {
	xing_ming_xue := htmlquery.FindOne(min_su_block, ".//text()[contains(., '姓名学：')]")

	if xing_ming_xue != nil {
		xing_ming_xue_str := strings.TrimSpace(htmlquery.InnerText(xing_ming_xue))

		if strings.Contains(xing_ming_xue_str, "非") {
			min_su.IsSurname = false
			if strings.Contains(xing_ming_xue_str, "男") {
				min_su.SurnameGender = "男"
			} else if strings.Contains(xing_ming_xue_str, "女") {
				min_su.SurnameGender = "女"
			} else {
				min_su.SurnameGender = "_"
			}
		} else {
			min_su.IsSurname = true
			if strings.Contains(xing_ming_xue_str, "男") {
				min_su.SurnameGender = "男"
			} else if strings.Contains(xing_ming_xue_str, "女") {
				min_su.SurnameGender = "女"
			} else {
				min_su.SurnameGender = "_"
			}
		}
	}

}
func parseLucky(wu_xing_ji_xiong_regular_str string) (ji_xiong string) {
	if strings.Contains(wu_xing_ji_xiong_regular_str, "寓意：") {
		ji_xiong = utf8string.NewString(strings.TrimSpace(strings.Split(wu_xing_ji_xiong_regular_str, "寓意：")[1])).Slice(0, 1)
	} else {
		ji_xiong = "_"
	}

	return ji_xiong
}
func parseWuXing(wu_xing_ji_xiong_regular_str string) (wu_xing string) {
	if strings.Contains(wu_xing_ji_xiong_regular_str, "五行：") {
		wu_xing = utf8string.NewString(strings.TrimSpace(strings.Split(wu_xing_ji_xiong_regular_str, "五行：")[1])).Slice(0, 1)
	} else {
		wu_xing = "_"
	}

	if wu_xing == "岁" {
		wu_xing = "水"
	}

	return wu_xing
}

//民俗 汉字五行： 吉凶寓意： 姓名学： 是否为常用字：
func parseDictInformation(exc *Excavator, unid int, html_node *html.Node) (err error) {
	min_su_block := htmlquery.FindOne(html_node, "//div[contains(@class, 'text16')]/span[contains(text(), '俗')]/..")

	if min_su_block == nil {
		panic("民俗参考区找不到")
	}

	wu_xing_ji_xiong_regular := htmlquery.FindOne(min_su_block, ".//text()[contains(., '常用字：')]")

	min_su := &models.MinSu{}

	if wu_xing_ji_xiong_regular != nil {
		wu_xing_ji_xiong_regular_str := string([]rune(strings.TrimSpace(htmlquery.InnerText(wu_xing_ji_xiong_regular))))

		wu_xing := parseWuXing(wu_xing_ji_xiong_regular_str)

		ji_xiong := parseLucky(wu_xing_ji_xiong_regular_str)

		chang_yong := parseRegular(wu_xing_ji_xiong_regular_str)

		min_su = &models.MinSu{
			WuXing: wu_xing,
			Lucky:  ji_xiong,
		}

		if chang_yong == "是" {
			min_su.Regular = true
		} else {
			min_su.Regular = false
		}
	}

	parseNameScience(min_su, min_su_block)

	if min_su != nil {
		err = InsertOrUpdate(exc.Db, min_su)
		if err != nil {
			panic(err)
		}

		_, err := exc.Db.Get(min_su)
		if err != nil {
			panic(err)
		}

		min_su_idr := models.MinSuId{
			Msid: min_su.Msid,
			Unid: unid,
		}

		err = InsertOrUpdate(exc.Db, &min_su_idr)
		if err != nil {
			panic(err)
		}
	}

	return nil
}

//基本解释
func parseComment(exc *Excavator, unid int, html_node *html.Node, ji_ben_block *html.Node) (err error) {
	ji_ben_empty := htmlquery.FindOne(ji_ben_block, "./font")

	hanCheng := &models.HanCheng{
		Unid: unid,
	}

	if ji_ben_empty == nil {
		ji_ben := htmlquery.Find(ji_ben_block, "./text()[not(normalize-space(.)='') and not(contains(., '笔画')) and not(contains(., '部首')) and not(contains(., '笔顺'))]")

		content_strs := []string{}

		for _, content := range ji_ben {
			content_line := strings.TrimSpace(htmlquery.InnerText(content))
			content_strs = append(content_strs, content_line)
		}

		//fix content mistake
		switch string(rune(unid)) {
		case "只":
			for idx, content_strs_line := range content_strs {
				if strings.Contains(content_strs_line, "（祇）") {
					content_strs[idx] = strings.ReplaceAll(content_strs_line, "（祇）", "（衹）")
				}
			}
		case "么":
			for idx, content_strs_line := range content_strs {
				if strings.Contains(content_strs_line, "（麽）") {
					content_strs[idx] = strings.ReplaceAll(content_strs_line, "（麽）", "（麼）")
				}
			}
		}

		if len(content_strs) > 0 {
			hanCheng.Content = content_strs

			err = InsertOrUpdate(exc.Db, hanCheng)
			if err != nil {
				panic(err)
			}
		}

	} else {
		return errors.New("基本解释为空")
	}

	return nil
}

//康熙字典解释
func parseDictComment(exc *Excavator, unid int, html_node *html.Node, touch bool) (kang_stroke int, err error) {

	kang_xi_dict := htmlquery.FindOne(html_node, "//div[@id='div_a4']/span[contains(text(), '康熙字典解释')]/..")

	if kang_xi_dict == nil {
		panic("康熙字典部分解析不到啦")
	}

	kang_xi_dict_brief := htmlquery.FindOne(kang_xi_dict, ".//strong")

	if kang_xi_dict_brief == nil {
		panic("康熙字典的brief解析不到")
	}

	kang_xi_dict_brief_str := strings.TrimSpace(htmlquery.InnerText(kang_xi_dict_brief))

	kang_xi_dict_brief_str_arr := strings.Split(kang_xi_dict_brief_str, "康熙笔画：")

	kang_xi_dict_brief_str_tail_arr := strings.Split(kang_xi_dict_brief_str_arr[1], "；")

	kang_xi_stroke_num, err := strconv.ParseUint(kang_xi_dict_brief_str_tail_arr[0], 10, bits.UintSize)
	if err != nil {
		panic(err)
	}

	kang_xi_stroke_num_int := int(kang_xi_stroke_num)

	if !touch {
		kang_xi_dict_brief_str_front_arr := strings.Split(kang_xi_dict_brief_str_arr[0], string(rune(unid)))

		kang_xi_dict_url := htmlquery.FindOne(kang_xi_dict_brief, ".//a")

		kang_xi_dict_content := htmlquery.Find(kang_xi_dict, ".//strong/following-sibling::text()[not(normalize-space(.)='')]")

		kang_xi_dict_content_strs := []string{}

		for _, kang_xi_dict_content_line := range kang_xi_dict_content {
			kang_xi_dict_content_strs = append(kang_xi_dict_content_strs, strings.TrimSpace(htmlquery.InnerText(kang_xi_dict_content_line)))
		}

		kang_xi := models.KangXi{
			Unid:    unid,
			Stroke:  kang_xi_stroke_num_int,
			Brief:   strings.TrimSpace(kang_xi_dict_brief_str_front_arr[0]),
			Content: kang_xi_dict_content_strs,
			Url:     htmlquery.SelectAttr(kang_xi_dict_url, "href"),
		}

		err = InsertOrUpdate(exc.Db, &kang_xi)
		if err != nil {
			panic(err)
		}
	}

	return kang_xi_stroke_num_int, nil
}
