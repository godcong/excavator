package excavator

import (
	"fmt"
	"strings"

	"github.com/antchfx/htmlquery"
	"github.com/godcong/excavator/models"
	"golang.org/x/net/html"
)

//常用编码
func parseBianMa(exc *Excavator, unid rune, html_node *html.Node) (err error) {

	bian_ma_block := htmlquery.FindOne(html_node, "//p[contains(@class, 'text16')]")

	if bian_ma_block == nil {
		panic("编码区找不到")
	}

	wu_bi_86 := htmlquery.FindOne(bian_ma_block, ".//span[contains(text(), '86')]/following-sibling::text()[1]")

	var wu_bi_86_str string

	if wu_bi_86 == nil {
		wu_bi_86_str = ""
	} else {
		wu_bi_86_str = strings.TrimSpace(htmlquery.InnerText(wu_bi_86))
	}

	wu_bi_98 := htmlquery.FindOne(bian_ma_block, ".//span[contains(text(), '98')]/following-sibling::text()[1]")

	var wu_bi_98_str string

	if wu_bi_98 == nil {
		wu_bi_98_str = ""
	} else {
		wu_bi_98_str = strings.TrimSpace(htmlquery.InnerText(wu_bi_98))
	}

	cang_jie := htmlquery.FindOne(bian_ma_block, ".//span[contains(text(), '仓颉')]/following-sibling::text()[1]")

	var cang_jie_str string

	if cang_jie == nil {
		cang_jie_str = ""
	} else {
		cang_jie_str = strings.TrimSpace(htmlquery.InnerText(cang_jie))
	}

	si_jiao := htmlquery.FindOne(bian_ma_block, ".//span[contains(text(), '四角')]/following-sibling::text()[1]")

	var si_jiao_str string

	if si_jiao == nil {
		si_jiao_str = ""
	} else {
		si_jiao_str = strings.TrimSpace(htmlquery.InnerText(si_jiao))
	}

	gui_fan := htmlquery.FindOne(bian_ma_block, ".//span[contains(text(), '规范汉字')]/following-sibling::text()[1]")

	var gui_fan_str string

	if gui_fan == nil {
		gui_fan_str = ""
	} else {
		gui_fan_str = strings.TrimSpace(htmlquery.InnerText(gui_fan))
	}

	bian_ma := models.BianMa{
		Unid:    unid,
		WuBi86:  wu_bi_86_str,
		WuBi98:  wu_bi_98_str,
		CangJie: cang_jie_str,
		SiJiao:  si_jiao_str,
		GuiFan:  gui_fan_str,
	}

	err = InsertOrUpdate(exc.Db, &bian_ma)
	if err != nil {
		panic(err)
	}

	return nil
}

//音韵
func parseYinYun(exc *Excavator, unid rune, html_node *html.Node, he_xin_block *html.Node) (err error) {
	yin_yun_block := htmlquery.FindOne(html_node, "//div[contains(@class, 'text16')]/span[contains(text(), '音韵')]/..")

	if yin_yun_block != nil {

		shang_gu := htmlquery.FindOne(yin_yun_block, ".//span[contains(text(), '上古')]/following-sibling::text()[1]")

		guang_yun := htmlquery.FindOne(yin_yun_block, ".//span[contains(text(), '广　韵')]/following-sibling::text()[1]")

		ping_shui := htmlquery.FindOne(yin_yun_block, ".//span[contains(text(), '平水')]/following-sibling::text()[1]")

		tang_yin := htmlquery.FindOne(yin_yun_block, ".//span[contains(text(), '唐　音')]/following-sibling::text()[1]")

		guo_yu := htmlquery.FindOne(yin_yun_block, ".//span[contains(text(), '国　语')]/following-sibling::text()[1]")

		yue_yu := htmlquery.Find(yin_yun_block, ".//span[contains(text(), '粤　语')]/following-sibling::script")

		min_nan := htmlquery.FindOne(yin_yun_block, ".//span[contains(text(), '闽南')]/following-sibling::text()[1]")

		pin_yin := htmlquery.Find(he_xin_block, ".//span[contains(@class, 'b') and contains(text(), '拼')]/following-sibling::span[1]/text()")

		zhu_yin := htmlquery.Find(he_xin_block, ".//span[contains(@class, 'b') and contains(text(), '注')]/following-sibling::span[1]/text()")

		var shang_gu_str string
		if shang_gu != nil {
			shang_gu_str = strings.TrimSpace(strings.Split(strings.TrimSpace(htmlquery.InnerText(shang_gu)), "]：")[1])
		}

		var guang_yun_str string
		if guang_yun != nil {
			guang_yun_str = strings.TrimSpace(strings.Split(strings.TrimSpace(htmlquery.InnerText(guang_yun)), "]：")[1])
		}

		var ping_shui_str string
		if ping_shui != nil {
			ping_shui_str = strings.TrimSpace(strings.Split(strings.TrimSpace(htmlquery.InnerText(ping_shui)), "]：")[1])
		}

		if len(shang_gu_str) > 0 || len(guang_yun_str) > 0 || len(ping_shui_str) > 0 {
			yin_yun := models.YinYun{
				Unid:     unid,
				ShangGu:  shang_gu_str,
				GuangYun: guang_yun_str,
				PingShui: ping_shui_str,
			}

			InsertOrUpdate(exc.Db, &yin_yun)
		}

		var tang_yin_str string
		if tang_yin != nil {
			tang_yin_str = strings.TrimSpace(strings.Split(strings.TrimSpace(htmlquery.InnerText(tang_yin)), "]：")[1])
		}

		var guo_yu_str string
		if guo_yu != nil {
			guo_yu_str = strings.TrimSpace(strings.Split(strings.TrimSpace(htmlquery.InnerText(guo_yu)), "]：")[1])
		}

		var min_nan_str string
		if min_nan != nil {
			min_nan_str = strings.TrimSpace(strings.Split(strings.TrimSpace(htmlquery.InnerText(min_nan)), "]：")[1])
		}

		tang_yin_strs := strings.Split(tang_yin_str, ",")

		guo_yu_strs := strings.Split(guo_yu_str, ",")

		yue_yu_strs := []string{}
		for idx, yue_yu_swf := range yue_yu {
			yue_yu_ele := htmlquery.FindOne(yue_yu_swf, ".//preceding-sibling::text()[1]")
			if idx == 0 {
				yue_yu_strs = append(yue_yu_strs, strings.TrimSpace(strings.Split(strings.TrimSpace(htmlquery.InnerText(yue_yu_ele)), "]：")[1]))
			} else {
				yue_yu_strs = append(yue_yu_strs, strings.TrimSpace(htmlquery.InnerText(yue_yu_ele)))
			}

		}

		min_nan_strs := strings.Split(min_nan_str, ",")

		pin_yin_strs := []string{}
		for _, pin_yin_ele := range pin_yin {
			pin_yin_ele_str := strings.TrimSpace(htmlquery.InnerText(pin_yin_ele))
			if strings.Contains(pin_yin_ele_str, "　") {
				for _, pin_yin_ele_str_sub := range strings.Split(pin_yin_ele_str, "　") {
					pin_yin_ele_str_sub_str := strings.TrimSpace(pin_yin_ele_str_sub)
					if len(pin_yin_ele_str_sub_str) > 0 {
						pin_yin_strs = append(pin_yin_strs, pin_yin_ele_str_sub_str)
					}
				}
			} else {
				pin_yin_strs = append(pin_yin_strs, pin_yin_ele_str)
			}
		}

		zhu_yin_strs := []string{}
		for _, zhu_yin_ele := range zhu_yin {
			zhu_yin_ele_str := strings.TrimSpace(htmlquery.InnerText(zhu_yin_ele))
			if strings.Contains(zhu_yin_ele_str, "　") {
				for _, zhu_yin_ele_str_sub := range strings.Split(zhu_yin_ele_str, "　") {
					zhu_yin_ele_str_sub_str := strings.TrimSpace(zhu_yin_ele_str_sub)
					if len(zhu_yin_ele_str_sub_str) > 0 {
						zhu_yin_strs = append(zhu_yin_strs, zhu_yin_ele_str_sub_str)
					}
				}
			} else {
				zhu_yin_strs = append(zhu_yin_strs, zhu_yin_ele_str)
			}
		}

		tang_yins := []*models.TangYin{}

		for _, tang_yin_strs_sub := range tang_yin_strs {
			if len(tang_yin_strs_sub) > 0 {
				tang_yins = append(tang_yins, &models.TangYin{
					TangYin: tang_yin_strs_sub,
				})
			}
		}

		for _, tang_yins_ele := range tang_yins {
			InsertOrUpdate(exc.Db, tang_yins_ele)
			exc.Db.Get(tang_yins_ele)
			tang_yin_id := models.TangYinId{
				Tid:  tang_yins_ele.Tid,
				Unid: unid,
			}
			InsertOrUpdate(exc.Db, &tang_yin_id)
		}

		guo_yus := []*models.PinYin{}

		for _, guo_yu_strs_sub := range guo_yu_strs {
			if len(guo_yu_strs_sub) > 0 {
				guo_yus = append(guo_yus, &models.PinYin{
					PinYin: guo_yu_strs_sub,
				})
			}
		}

		for _, guo_yus_ele := range guo_yus {
			InsertOrUpdate(exc.Db, guo_yus_ele)
			exc.Db.Get(guo_yus_ele)
			guo_yu_id := models.GuoYuId{
				Pid:  guo_yus_ele.Pid,
				Unid: unid,
			}
			InsertOrUpdate(exc.Db, &guo_yu_id)
		}

		pin_yins := []*models.PinYin{}
		for _, pin_yin_strs_sub := range pin_yin_strs {
			pin_yins = append(pin_yins, &models.PinYin{
				PinYin: pin_yin_strs_sub,
			})
		}

		for _, pin_yin_ele := range pin_yins {
			InsertOrUpdate(exc.Db, pin_yin_ele)
			exc.Db.Get(pin_yin_ele)
			pin_yin_id := models.PinYinId{
				Pid:  pin_yin_ele.Pid,
				Unid: unid,
			}
			InsertOrUpdate(exc.Db, &pin_yin_id)
		}

		zhu_yins := []*models.ZhuYin{}

		for _, zhu_yin_strs_sub := range zhu_yin_strs {
			zhu_yins = append(zhu_yins, &models.ZhuYin{
				ZhuYin: zhu_yin_strs_sub,
			})
		}

		for _, zhu_yin_ele := range zhu_yins {
			InsertOrUpdate(exc.Db, zhu_yin_ele)
			exc.Db.Get(zhu_yin_ele)
			zhu_yin_id := models.ZhuYinId{
				Zid:  zhu_yin_ele.Zid,
				Unid: unid,
			}
			InsertOrUpdate(exc.Db, &zhu_yin_id)
		}

		yue_yus := []*models.YueYin{}

		for _, yue_yu_strs_sub := range yue_yu_strs {
			yue_yus = append(yue_yus, &models.YueYin{
				YueYin: yue_yu_strs_sub,
			})
		}

		for _, yue_yu_ele := range yue_yus {
			InsertOrUpdate(exc.Db, yue_yu_ele)
			exc.Db.Get(yue_yu_ele)
			yue_yu_id := models.YueYinId{
				Yid:  yue_yu_ele.Yid,
				Unid: unid,
			}
			InsertOrUpdate(exc.Db, &yue_yu_id)
		}

		min_nans := []*models.MinNanYin{}

		for _, min_nan_strs_sub := range min_nan_strs {
			if len(min_nan_strs_sub) > 0 {
				min_nans = append(min_nans, &models.MinNanYin{
					MinNanYin: min_nan_strs_sub,
				})
			}
		}

		for _, min_nan_ele := range min_nans {
			InsertOrUpdate(exc.Db, min_nan_ele)
			exc.Db.Get(min_nan_ele)
			min_nan_id := models.MinNanYinId{
				Mnid: min_nan_ele.Mnid,
				Unid: unid,
			}
			InsertOrUpdate(exc.Db, &min_nan_id)
		}
	}

	return nil
}

//索引参考
func parseSuoYin(exc *Excavator, unid rune, html_node *html.Node) (err error) {

	suo_yin_block := htmlquery.FindOne(html_node, "//div[contains(@class, 'text16')]/span[contains(text(), '索引')]/..")

	if suo_yin_block != nil {

		gu_wen := htmlquery.FindOne(suo_yin_block, ".//span[contains(text(), '古文')]/following-sibling::text()[1]")

		gu_xun := htmlquery.FindOne(suo_yin_block, ".//span[contains(text(), '故训')]/following-sibling::text()[1]")

		shuo_wen_suo_yin := htmlquery.FindOne(suo_yin_block, ".//span[contains(text(), '说文')]/following-sibling::text()[1]")

		kang_xi_suo_yin := htmlquery.FindOne(suo_yin_block, ".//span[contains(text(), '康熙')]/following-sibling::text()[1]")

		han_yu_suo_yin := htmlquery.FindOne(suo_yin_block, ".//span[contains(text(), '汉语')]/following-sibling::text()[1]")

		ci_hai := htmlquery.FindOne(suo_yin_block, ".//span[contains(text(), '辞　海')]/following-sibling::text()[1]")

		var gu_wen_str string
		if gu_wen != nil {
			gu_wen_str = strings.TrimSpace(strings.Split(strings.TrimSpace(htmlquery.InnerText(gu_wen)), "]：")[1])
		}

		var gu_xun_str string
		if gu_xun != nil {
			gu_xun_str = strings.TrimSpace(strings.Split(strings.TrimSpace(htmlquery.InnerText(gu_xun)), "]：")[1])
		}

		var shuo_wen_suo_yin_str string
		if shuo_wen_suo_yin != nil {
			shuo_wen_suo_yin_str = strings.TrimSpace(strings.Split(strings.TrimSpace(htmlquery.InnerText(shuo_wen_suo_yin)), "]：")[1])
		}

		var kang_xi_suo_yin_str string
		if kang_xi_suo_yin != nil {
			kang_xi_suo_yin_str = strings.TrimSpace(strings.Split(strings.TrimSpace(htmlquery.InnerText(kang_xi_suo_yin)), "]：")[1])
		}

		var han_yu_suo_yin_str string
		if han_yu_suo_yin != nil {
			han_yu_suo_yin_str = strings.TrimSpace(strings.Split(strings.TrimSpace(htmlquery.InnerText(han_yu_suo_yin)), "]：")[1])
		}

		var ci_hai_str string
		if ci_hai != nil {
			ci_hai_str = strings.TrimSpace(strings.Split(strings.TrimSpace(htmlquery.InnerText(ci_hai)), "]：")[1])
		}

		if len(gu_wen_str) > 0 || len(gu_xun_str) > 0 || len(shuo_wen_suo_yin_str) > 0 || len(kang_xi_suo_yin_str) > 0 || len(han_yu_suo_yin_str) > 0 || len(ci_hai_str) > 0 {
			suo_yin := models.SuoYin{
				Unid: unid,
				Gu:   gu_wen_str,
				Xun:  gu_xun_str,
				Shuo: shuo_wen_suo_yin_str,
				Kang: kang_xi_suo_yin_str,
				Han:  han_yu_suo_yin_str,
				Ci:   ci_hai_str,
			}

			InsertOrUpdate(exc.Db, &suo_yin)
		}
	}

	return nil
}

//详细解释，新华字典
func parseXiangXi(exc *Excavator, unid rune, html_node *html.Node) (err error) {

	xin_hua_block := htmlquery.FindOne(html_node, "//div[@id='div_a2']")

	xin_hua_empty := htmlquery.Find(xin_hua_block, ".//font")

	if xin_hua_empty == nil {
		xin_hua_lines := htmlquery.Find(xin_hua_block, ".//span/following-sibling::text()[not(normalize-space(.)='')]")

		xin_hua_strs := []string{}

		for _, xin_hua_line := range xin_hua_lines {
			xin_hua_strs = append(xin_hua_strs, strings.TrimSpace(htmlquery.InnerText(xin_hua_line)))
		}

		xin_hua := models.XinHua{
			Unid:    unid,
			Content: xin_hua_strs,
		}

		InsertOrUpdate(exc.Db, &xin_hua)
	}

	return nil
}

//汉语字典
func parseHanYuDaZiDian(exc *Excavator, unid rune, html_node *html.Node) (err error) {
	han_yu_block := htmlquery.FindOne(html_node, "//div[@id='div_a3']")

	han_yu_empty := htmlquery.FindOne(han_yu_block, ".//font")

	if han_yu_empty == nil {
		han_yu_strs := []string{}

		for _, han_yu_line := range htmlquery.Find(han_yu_block, ".//p/text()") {
			han_yu_line_str := strings.TrimSpace(htmlquery.InnerText(han_yu_line))

			han_yu_strs = append(han_yu_strs, han_yu_line_str)
		}

		han_yu := models.HanDa{
			Unid:    unid,
			Content: han_yu_strs,
		}

		InsertOrUpdate(exc.Db, &han_yu)
	}

	return nil
}

//说文解字
func parseShuoWenJieZi(exc *Excavator, unid rune, html_node *html.Node, shuo_wen_block *html.Node) (err error) {
	shuo_wen_font_new := htmlquery.FindOne(shuo_wen_block, ".//span[boolean(@style)]")

	shuo_wen_font_old := htmlquery.FindOne(shuo_wen_block, ".//img")

	var shuo_wen_font_old_str string

	if shuo_wen_font_old != nil {
		shuo_wen_font_old_str = htmlquery.SelectAttr(shuo_wen_font_old, "src")
	}

	shuo_wen_empty := htmlquery.FindOne(shuo_wen_block, ".//font")

	if shuo_wen_empty == nil {

		var shuo_wen_brief_str string

		if shuo_wen_font_new != nil {
			shuo_wen_brief := htmlquery.FindOne(shuo_wen_block, ".//span[boolean(@style)]/following-sibling::text()[1]")

			if shuo_wen_brief != nil {
				shuo_wen_brief_str = strings.TrimSpace(htmlquery.InnerText(shuo_wen_brief))

				shuo_wen_lines := htmlquery.Find(shuo_wen_block, ".//span[boolean(@style)]/following-sibling::text()[1]/following-sibling::text()[not(normalize-space(.)='') and not(contains(., '古文"+string(rune(unid))+"'))]")

				shuo_wen_content_strs := []string{}

				for _, shuo_wen_line := range shuo_wen_lines {
					if shuo_wen_line.Type == html.TextNode {
						shuo_wen_content_str := strings.TrimSpace(htmlquery.InnerText(shuo_wen_line))

						shuo_wen_content_strs = append(shuo_wen_content_strs, shuo_wen_content_str)
					}
				}

				shuo_wen := models.ShuoWen{
					Unid:    unid,
					Brief:   shuo_wen_brief_str,
					Content: shuo_wen_content_strs,
					Url:     shuo_wen_font_old_str,
				}

				InsertOrUpdate(exc.Db, &shuo_wen)
			}
		}

	}

	return nil
}

//+ 字源演变
func parseYanBian(exc *Excavator, unid rune, html_node *html.Node, shuo_wen_block *html.Node) (err error) {
	jia_gu_pic := htmlquery.FindOne(shuo_wen_block, ".//li/text()[contains(., '甲骨')]/../p")

	jin_wen_pic := htmlquery.FindOne(shuo_wen_block, ".//li/text()[contains(., '金文')]/../p")

	xiao_zhuan_pic := htmlquery.FindOne(shuo_wen_block, ".//li/text()[contains(., '小篆')]/../p")

	kai_ti_pic := htmlquery.FindOne(shuo_wen_block, ".//li/text()[contains(., '楷体')]/../p")

	var jia_gu_url string

	if jia_gu_pic != nil {
		jia_gu_url = strings.Split(strings.Split(htmlquery.SelectAttr(jia_gu_pic, "style"), "background:url(")[1], ")")[0]
	}

	var jin_wen_url string

	if jin_wen_pic != nil {
		jin_wen_url = strings.Split(strings.Split(htmlquery.SelectAttr(jin_wen_pic, "style"), "background:url(")[1], ")")[0]
	}

	var xiao_zhuan_url string

	if xiao_zhuan_pic != nil {
		xiao_zhuan_url = strings.Split(strings.Split(htmlquery.SelectAttr(xiao_zhuan_pic, "style"), "background:url(")[1], ")")[0]
	}

	var kai_ti_url string

	if kai_ti_pic != nil {
		kai_ti_url = strings.Split(strings.Split(htmlquery.SelectAttr(kai_ti_pic, "style"), "background:url(")[1], ")")[0]
	}

	if len(jia_gu_url) > 0 || len(jin_wen_url) > 0 || len(xiao_zhuan_url) > 0 || len(kai_ti_url) > 0 {
		yan_bian := models.YanBian{
			Unid:         unid,
			JiaGuUrl:     jia_gu_url,
			JinUrl:       jin_wen_url,
			XiaoZhuanUrl: xiao_zhuan_url,
			KaiUrl:       kai_ti_url,
		}

		InsertOrUpdate(exc.Db, &yan_bian)
	}

	return nil
}

//成语
func parseChengYu(exc *Excavator, unid rune, html_node *html.Node, cis map[string]bool) (err error) {

	cheng_yu_links := htmlquery.Find(html_node, "//ul[contains(@class, 'chengyu')]//a")

	for _, cheng_yu_link := range cheng_yu_links {
		url_str := htmlquery.SelectAttr(cheng_yu_link, "href")
		cheng_yu_front := htmlquery.FindOne(cheng_yu_link, ".//font[1]/preceding-sibling::text()[not(normalize-space(.)='')][1]")
		cheng_yu_mid_list := htmlquery.Find(cheng_yu_link, ".//font")

		var cheng_yu_str string
		if cheng_yu_front != nil {
			cheng_yu_str = strings.TrimSpace(htmlquery.InnerText(cheng_yu_front))
		}

		for idx, cheng_yu_mid := range cheng_yu_mid_list {
			cheng_yu_mid_next_text := htmlquery.FindOne(cheng_yu_link, fmt.Sprintf(".//font[%d]/following-sibling::text()[not(normalize-space(.)='')][1]", idx+1))
			cheng_yu_mid_next_text_font := htmlquery.Find(cheng_yu_link, fmt.Sprintf(".//font[%d]/following-sibling::text()[not(normalize-space(.)='')][1]/following-sibling::font", idx+1))
			cheng_yu_mid_next_font := htmlquery.Find(cheng_yu_link, fmt.Sprintf(".//font[%d]/following-sibling::font", idx+1))

			if cheng_yu_mid_next_font == nil && cheng_yu_mid_next_text_font == nil {
				if cheng_yu_mid_next_text != nil {
					cheng_yu_str += strings.TrimSpace(htmlquery.InnerText(cheng_yu_mid)) + strings.TrimSpace(htmlquery.InnerText(cheng_yu_mid_next_text))
				} else {
					cheng_yu_str += strings.TrimSpace(htmlquery.InnerText(cheng_yu_mid))
				}
			} else if cheng_yu_mid_next_font != nil && cheng_yu_mid_next_text_font == nil {
				cheng_yu_str += strings.TrimSpace(htmlquery.InnerText(cheng_yu_mid))
			} else {
				if len(cheng_yu_mid_next_font) == len(cheng_yu_mid_next_text_font) {
					cheng_yu_str += strings.TrimSpace(htmlquery.InnerText(cheng_yu_mid)) + strings.TrimSpace(htmlquery.InnerText(cheng_yu_mid_next_text))
				} else {
					cheng_yu_str += strings.TrimSpace(htmlquery.InnerText(cheng_yu_mid))
				}
			}
		}

		if len(cheng_yu_str) != 0 {

			cheng_yu := models.ChengYu{
				ChengYu: cheng_yu_str,
				Url:     url_str,
			}

			if !cis[cheng_yu.ChengYu] {
				InsertOrUpdate(exc.Db, &cheng_yu)

				exc.Db.Get(&cheng_yu)

				cheng_yu_id := models.ChengYuId{
					Cyid: cheng_yu.Cyid,
					Unid: unid,
				}

				InsertOrUpdate(exc.Db, &cheng_yu_id)

				cis[cheng_yu.ChengYu] = true
			}
		}

	}

	return nil
}

//诗词
func parseShiCi(exc *Excavator, unid rune, html_node *html.Node, cis map[string]bool) (err error) {
	shi_ci_links := htmlquery.Find(html_node, "//ul[contains(@class, 'shici')]//a")

	for _, shi_ci_link := range shi_ci_links {
		url_str := htmlquery.SelectAttr(shi_ci_link, "href")
		shi_ci_front := htmlquery.FindOne(shi_ci_link, ".//font[1]/preceding-sibling::text()[not(normalize-space(.)='')][1]")
		shi_ci_mid_list := htmlquery.Find(shi_ci_link, ".//font")

		var shi_ci_str string
		if shi_ci_front != nil {
			shi_ci_str = strings.ReplaceAll(strings.TrimSpace(htmlquery.InnerText(shi_ci_front)), " ", "")
			shi_ci_str = strings.ReplaceAll(shi_ci_str, "　", "")
		}

		for idx, shi_ci_mid := range shi_ci_mid_list {
			shi_ci_mid_next_text := htmlquery.FindOne(shi_ci_link, fmt.Sprintf(".//font[%d]/following-sibling::text()[not(normalize-space(.)='')][1]", idx+1))
			shi_ci_mid_next_text_font := htmlquery.Find(shi_ci_link, fmt.Sprintf(".//font[%d]/following-sibling::text()[not(normalize-space(.)='')][1]/following-sibling::font", idx+1))
			shi_ci_mid_next_font := htmlquery.Find(shi_ci_link, fmt.Sprintf(".//font[%d]/following-sibling::font", idx+1))

			if shi_ci_mid_next_font == nil && shi_ci_mid_next_text_font == nil {
				if shi_ci_mid_next_text != nil {
					shi_ci_mid_next_text_str := strings.ReplaceAll(strings.TrimSpace(htmlquery.InnerText(shi_ci_mid_next_text)), " ", "")
					shi_ci_mid_next_text_str = strings.ReplaceAll(shi_ci_mid_next_text_str, "　", "")
					shi_ci_str += strings.TrimSpace(htmlquery.InnerText(shi_ci_mid)) + shi_ci_mid_next_text_str
				} else {
					shi_ci_str += strings.TrimSpace(htmlquery.InnerText(shi_ci_mid))
				}
			} else if shi_ci_mid_next_font != nil && shi_ci_mid_next_text_font == nil {
				shi_ci_str += strings.TrimSpace(htmlquery.InnerText(shi_ci_mid))
			} else {
				if len(shi_ci_mid_next_font) == len(shi_ci_mid_next_text_font) {
					shi_ci_mid_next_text_str := strings.ReplaceAll(strings.TrimSpace(htmlquery.InnerText(shi_ci_mid_next_text)), " ", "")
					shi_ci_mid_next_text_str = strings.ReplaceAll(shi_ci_mid_next_text_str, "　", "")
					shi_ci_str += strings.TrimSpace(htmlquery.InnerText(shi_ci_mid)) + shi_ci_mid_next_text_str
				} else {
					shi_ci_str += strings.TrimSpace(htmlquery.InnerText(shi_ci_mid))
				}
			}
		}

		if len(shi_ci_str) != 0 {
			shi_ci := models.ShiCi{
				ShiCi: shi_ci_str,
				Url:   url_str,
			}

			if !cis[shi_ci.ShiCi] {
				InsertOrUpdate(exc.Db, &shi_ci)

				exc.Db.Get(&shi_ci)

				shi_ci_id := models.ShiCiId{
					Scid: shi_ci.Scid,
					Unid: unid,
				}

				InsertOrUpdate(exc.Db, &shi_ci_id)

				cis[shi_ci.ShiCi] = true
			}

		}

	}

	return nil
}

//词语
func parseCiYu(exc *Excavator, unid rune, html_node *html.Node, cis map[string]bool) (err error) {
	ci_yu_links := htmlquery.Find(html_node, "//ul[contains(@class, 'ciyu')]//a")

	for _, ci_yu_link := range ci_yu_links {
		url_str := htmlquery.SelectAttr(ci_yu_link, "href")
		ci_yu_front := htmlquery.FindOne(ci_yu_link, ".//font[1]/preceding-sibling::text()[not(normalize-space(.)='')][1]")
		ci_yu_mid_list := htmlquery.Find(ci_yu_link, ".//font")

		var ci_yu_str string
		if ci_yu_front != nil {
			ci_yu_str = strings.TrimSpace(htmlquery.InnerText(ci_yu_front))
		}

		for idx, ci_yu_mid := range ci_yu_mid_list {
			ci_yu_mid_next_text := htmlquery.FindOne(ci_yu_link, fmt.Sprintf(".//font[%d]/following-sibling::text()[not(normalize-space(.)='')][1]", idx+1))
			ci_yu_mid_next_text_font := htmlquery.Find(ci_yu_link, fmt.Sprintf(".//font[%d]/following-sibling::text()[not(normalize-space(.)='')][1]/following-sibling::font", idx+1))
			ci_yu_mid_next_font := htmlquery.Find(ci_yu_link, fmt.Sprintf(".//font[%d]/following-sibling::font", idx+1))

			if ci_yu_mid_next_font == nil && ci_yu_mid_next_text_font == nil {
				if ci_yu_mid_next_text != nil {
					ci_yu_str += strings.TrimSpace(htmlquery.InnerText(ci_yu_mid)) + strings.TrimSpace(htmlquery.InnerText(ci_yu_mid_next_text))
				} else {
					ci_yu_str += strings.TrimSpace(htmlquery.InnerText(ci_yu_mid))
				}
			} else if ci_yu_mid_next_font != nil && ci_yu_mid_next_text_font == nil {
				ci_yu_str += strings.TrimSpace(htmlquery.InnerText(ci_yu_mid))
			} else {
				if len(ci_yu_mid_next_font) == len(ci_yu_mid_next_text_font) {
					ci_yu_str += strings.TrimSpace(htmlquery.InnerText(ci_yu_mid)) + strings.TrimSpace(htmlquery.InnerText(ci_yu_mid_next_text))
				} else {
					ci_yu_str += strings.TrimSpace(htmlquery.InnerText(ci_yu_mid))
				}
			}
		}

		if len(ci_yu_str) != 0 {
			ci_yu := models.CiYu{
				CiYu: ci_yu_str,
				Url:  url_str,
			}

			if !cis[ci_yu.CiYu] {
				InsertOrUpdate(exc.Db, &ci_yu)

				exc.Db.Get(&ci_yu)

				ci_yu_id := models.CiYuId{
					Cid:  ci_yu.Cid,
					Unid: unid,
				}

				InsertOrUpdate(exc.Db, &ci_yu_id)

				cis[ci_yu.CiYu] = true
			}
		}

	}

	return nil
}
