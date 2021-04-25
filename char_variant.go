package excavator

import (
	"errors"
	"excavator/models"
	"fmt"
	"math/bits"
	"strconv"
	"strings"

	"github.com/antchfx/htmlquery"
	"github.com/godcong/fate"
	"golang.org/x/net/html"
	"xorm.io/xorm"
)

//基本解释中的异体字信息
func parseJiBenVariant(exc *Excavator, unid int, html_node *html.Node, ji_ben_block *html.Node) (err error) {
	variant_tongs := htmlquery.Find(ji_ben_block, "./text()[(contains(., '同“') and not(contains(., '古同“'))) or contains(., '见“')]")

	for _, variant_tong := range variant_tongs {
		v_tong_str := strings.TrimSpace(htmlquery.InnerText(variant_tong))

		v_char_str := strings.Split(strings.Split(v_tong_str, "“")[1], "”")[0]

		if len(v_char_str) == 0 {
			variant_tong_break := htmlquery.FindOne(ji_ben_block, "./text()[(contains(., '同“') and not(contains(., '古同“'))) or contains(., '见“')]/following-sibling::text()[not(normalize-space(.)='')][1]")

			if variant_tong_break == nil {
				return errors.New("基本解释同义字格式不对")
			}

			variant_tong_break_str := htmlquery.InnerText(variant_tong_break)

			v_char_str = strings.TrimSpace(strings.Split(variant_tong_break_str, "”")[0])

			if len(v_char_str) == 0 {
				return errors.New("基本解释同义字格式不对")
			}
		}

		v_char_str_unicode := int(([]rune(v_char_str))[0])

		v_char_str_unicode_str := strings.ToUpper(strconv.FormatUint(uint64(v_char_str_unicode), 16))

		v_uc := models.UnihanChar{
			Unid:       v_char_str_unicode,
			UnicodeHex: v_char_str_unicode_str,
		}

		err = InsertOrUpdate(exc.Db, &v_uc)
		if err != nil {
			panic(err)
		}

		v_idr := models.VariantId{
			Unid:  unid,
			UnidS: v_char_str_unicode,
		}

		err = InsertOrUpdate(exc.Db, &v_idr)
		if err != nil {
			panic(err)
		}
	}

	variant_gus := htmlquery.Find(ji_ben_block, "./text()[contains(., '古同“')]")

	for _, variant_gu := range variant_gus {
		v_gu := models.VariantGu{
			Unid: unid,
		}

		err = InsertOrUpdate(exc.Db, &v_gu)
		if err != nil {
			panic(err)
		}

		v_gu_str := strings.TrimSpace(htmlquery.InnerText(variant_gu))

		v_char_str := strings.Split(strings.Split(v_gu_str, "“")[1], "”")[0]

		if len(v_char_str) == 0 {
			variant_gu_break := htmlquery.FindOne(ji_ben_block, "./text()[contains(., '古同“')]/following-sibling::text()[not(normalize-space(.)='')][1]")

			if variant_gu_break == nil {
				panic("基本解释古同义字格式不对")
			}

			variant_gu_break_str := htmlquery.InnerText(variant_gu_break)

			v_char_str = strings.TrimSpace(strings.Split(variant_gu_break_str, "”")[0])
		}

		v_char_str_unicode := int(([]rune(v_char_str))[0])

		v_char_str_unicode_str := strings.ToUpper(strconv.FormatUint(uint64(v_char_str_unicode), 16))

		v_uc := models.UnihanChar{
			Unid:       v_char_str_unicode,
			UnicodeHex: v_char_str_unicode_str,
		}

		err = InsertOrUpdate(exc.Db, &v_uc)
		if err != nil {
			panic(err)
		}

		v_idr := models.VariantId{
			Unid:  unid,
			UnidS: v_char_str_unicode,
		}

		err = InsertOrUpdate(exc.Db, &v_idr)
		if err != nil {
			panic(err)
		}
	}

	return nil
}

//简体字，繁体字集 和 异体字集
func parseVariant(exc *Excavator, unid int, html_node *html.Node, he_xin_block *html.Node, ji_ben_block *html.Node) (err error) {
	jian_ti_list := htmlquery.Find(he_xin_block, ".//span[contains(@class, 'b') and contains(text(), '简体字：')]/following-sibling::a")

	jian_yi_ti_list := htmlquery.Find(he_xin_block, ".//span[contains(@class, 'b') and contains(text(), '简体字：')]/following-sibling::span[contains(@class, 'b') and contains(text(), '异')]/following-sibling::a")

	gu_list := map[int]*html.Node{}

	if jian_ti_list != nil {
		jian_ti := jian_ti_list[0]
		jian_ti_str := strings.TrimSpace(htmlquery.InnerText(jian_ti))
		j_unid := int(([]rune(jian_ti_str))[0])

		if jian_yi_ti_list != nil {
			if len(jian_ti_list)-len(jian_yi_ti_list) > 1 {
				panic("有多个简体字")
			}
			for _, yi_ti := range jian_yi_ti_list {
				yi_ti_str := strings.TrimSpace(htmlquery.InnerText(yi_ti))

				v_unid := int(([]rune(yi_ti_str))[0])

				if v_unid == j_unid {
					continue
				}

				gu_list[v_unid] = yi_ti
			}
		}

		uc := models.UnihanChar{
			Unid:       j_unid,
			UnicodeHex: strings.ToUpper(strconv.FormatUint(uint64(j_unid), 16)),
		}

		err = InsertOrUpdate(exc.Db, &uc)
		if err != nil {
			panic(err)
		}

		hcc := models.HanChengChar{
			Unid: j_unid,
			Url:  htmlquery.SelectAttr(jian_ti, "href"),
		}

		err = InsertOrUpdate(exc.Db, &hcc)
		if err != nil {
			panic(err)
		}

		v_idr := models.VariantId{
			Unid:  unid,
			UnidS: j_unid,
		}

		err = InsertOrUpdate(exc.Db, &v_idr)
		if err != nil {
			panic(err)
		}
	}

	fan_list := map[int]*html.Node{}

	fan_ti_list := htmlquery.Find(he_xin_block, ".//span[contains(@class, 'b') and contains(text(), '繁体字：')]/following-sibling::a")

	fan_yi_ti_list := htmlquery.Find(he_xin_block, ".//span[contains(@class, 'b') and contains(text(), '繁体字：')]/following-sibling::span[contains(@class, 'b') and contains(text(), '异')]/following-sibling::a")

	if fan_ti_list != nil {
		//有繁体字

		if fan_yi_ti_list != nil {
			//繁体字，异体字

			fan_ti_list = htmlquery.Find(he_xin_block, ".//span[contains(@class, 'b') and contains(text(), '异')]/preceding-sibling::a")

			for _, yi_ti := range fan_yi_ti_list {
				yi_ti_str := strings.TrimSpace(htmlquery.InnerText(yi_ti))

				v_unid := int(([]rune(yi_ti_str))[0])

				var is_fan_ti bool = false
				for _, fan_ti := range fan_ti_list {
					fan_ti_str := strings.TrimSpace(htmlquery.InnerText(fan_ti))
					f_unid := int(([]rune(fan_ti_str))[0])

					if v_unid == f_unid {
						is_fan_ti = true
						break
					}
				}

				if !is_fan_ti {
					gu_list[v_unid] = yi_ti
				}

			}

		}

		for _, yi_ti := range fan_ti_list {
			yi_ti_str := strings.TrimSpace(htmlquery.InnerText(yi_ti))

			v_unid := int(([]rune(yi_ti_str))[0])

			fan_list[v_unid] = yi_ti
		}

	}

	if jian_ti_list == nil && fan_ti_list == nil {
		yi_ti_list := htmlquery.Find(he_xin_block, ".//span[contains(@class, 'b') and contains(text(), '异')]/following-sibling::a")
		//无繁体字，但有异体字
		for _, yi_ti := range yi_ti_list {
			yi_ti_str := strings.TrimSpace(htmlquery.InnerText(yi_ti))

			v_unid := int(([]rune(yi_ti_str))[0])

			gu_list[v_unid] = yi_ti
		}
	}

	yi_list := map[int]*html.Node{}

	for uid, gu_ti := range gu_list {
		yi_list[uid] = gu_ti

		gu := models.VariantGu{
			Unid: uid,
		}

		err = InsertOrUpdate(exc.Db, &gu)
		if err != nil {
			panic(err)
		}
	}

	for uid, fan_ti := range fan_list {
		yi_list[uid] = fan_ti
	}

	for v_unid, yi_ti := range yi_list {
		v_uc := models.UnihanChar{
			Unid:       v_unid,
			UnicodeHex: strings.ToUpper(strconv.FormatUint(uint64(v_unid), 16)),
		}

		err = InsertOrUpdate(exc.Db, &v_uc)
		if err != nil {
			panic(err)
		}

		v_hcc := models.HanChengChar{
			Unid: v_unid,
			Url:  htmlquery.SelectAttr(yi_ti, "href"),
		}

		err = InsertOrUpdate(exc.Db, &v_hcc)
		if err != nil {
			panic(err)
		}

		v_idr := models.VariantId{
			Unid:  v_unid,
			UnidS: unid,
		}

		err = InsertOrUpdate(exc.Db, &v_idr)
		if err != nil {
			panic(err)
		}
	}

	return nil
}

//从康熙笔画列表中解析出姓名学笔画对应的字和笔画
//如果结果有多个字，那么解析结果依赖于页面中康熙字典的对应字段
func parseKangXiStroke(unid int, he_xin_block *html.Node) (kang map[rune]int) {
	kang_xi_stroke := htmlquery.FindOne(he_xin_block, ".//span[contains(@class, 'b') and contains(text(), '康')]/following-sibling::text()[1]")

	if kang_xi_stroke != nil {
		kang_xi_stroke_str := strings.TrimSpace(htmlquery.InnerText(kang_xi_stroke))

		kang_xi_stroke_str_tail := strings.ReplaceAll(kang_xi_stroke_str, "(", "")
		kang_xi_stroke_str_mid := strings.ReplaceAll(kang_xi_stroke_str_tail, ")", "")
		kang_xi_stroke_str_arr := strings.Split(kang_xi_stroke_str_mid, "；")

		fan_str_list := map[rune]int{}

		is_kangxi := false
		kang_stroke := 0

		for _, kang_xi_stroke_str_ele := range kang_xi_stroke_str_arr {
			kang_xi_stroke_str_ele_str := strings.TrimSpace(kang_xi_stroke_str_ele)
			if len(kang_xi_stroke_str_ele_str) > 0 {
				kang_xi_stroke_str_ele_arr := strings.Split(kang_xi_stroke_str_ele_str, ":")
				kang_xi_stroke_str_ele_ch := strings.TrimSpace(kang_xi_stroke_str_ele_arr[0])
				kang_xi_stroke_str_ele_num := strings.TrimSpace(kang_xi_stroke_str_ele_arr[1])

				kang_xi_stroke_str_ele_ch_rune := ([]rune(kang_xi_stroke_str_ele_ch))[0]

				kang_xi_stroke_str_ele_num_int, err := strconv.ParseUint(kang_xi_stroke_str_ele_num, 10, bits.UintSize)
				if err != nil {
					panic(err)
				}

				if rune(unid) == kang_xi_stroke_str_ele_ch_rune {
					is_kangxi = true
					kang_stroke = int(kang_xi_stroke_str_ele_num_int)
					continue
				}

				fan_str_list[kang_xi_stroke_str_ele_ch_rune] = int(kang_xi_stroke_str_ele_num_int)
			}
		}

		kang = map[rune]int{}

		if len(fan_str_list) > 1 {
			fmt.Println("有多个康熙繁体字")
			kang = fan_str_list
		} else if len(fan_str_list) == 1 {
			for my_rune, my_stroke := range fan_str_list {
				kang[my_rune] = my_stroke
			}
		} else if len(fan_str_list) == 0 {
			if is_kangxi {
				kang[rune(unid)] = kang_stroke
			}
		}
	}

	return kang
}

//整理变体序列表
func variantChars(exc *Excavator) (err error) {
	v_id_adds := []*models.VariantId{}
	v_id_dels := []*models.VariantId{}
	v_id_tmp := new(models.VariantId)
	v_id_rows, err := exc.Db.Rows(v_id_tmp)
	if err != nil {
		panic(err)
	}
	defer v_id_rows.Close()
	for v_id_rows.Next() {

		v_id := models.VariantId{}

		err = v_id_rows.Scan(&v_id)
		if err != nil {
			panic(err)
		}

		v_id_s_gots := []models.VariantId{}

		err := exc.Db.Where("Unid = ?", v_id.UnidS).Find(&v_id_s_gots)
		if err != nil {
			panic(err)
		}

		for _, v_id_s_got := range v_id_s_gots {
			v_id_adds = append(v_id_adds, &models.VariantId{
				Rid:   v_id.Rid,
				Unid:  v_id.Unid,
				UnidS: v_id_s_got.UnidS,
			})
		}

		if len(v_id_s_gots) > 0 {
			v_id_dels = append(v_id_dels, &v_id)
		}
	}

	for _, v_id_add := range v_id_adds {
		InsertOrUpdate(exc.Db, v_id_add)
	}

	for _, v_id_del := range v_id_dels {
		_, err := exc.Db.Delete(v_id_del)
		if err != nil {
			panic(err)
		}
	}

	return nil
}

//整理最简字
func simplifyChars(exc *Excavator) (err error) {
	han_cheng_char_tmp := new(models.HanChengChar)
	han_cheng_char_rows, err := exc.Db.Rows(han_cheng_char_tmp)
	if err != nil {
		panic(err)
	}
	han_chars := []*models.HanChar{}
	my_chars := []*fate.Character{}
	defer han_cheng_char_rows.Close()
	for han_cheng_char_rows.Next() {
		err = han_cheng_char_rows.Scan(han_cheng_char_tmp)
		if err != nil {
			panic(err)
		}

		v_id_tmp := models.VariantId{
			Unid: han_cheng_char_tmp.Unid,
		}

		has, err := exc.Db.Get(&v_id_tmp)
		if err != nil {
			panic(err)
		}

		if !has {
			han_char := models.HanChar{
				Unid: han_cheng_char_tmp.Unid,
				Ch:   string(rune(han_cheng_char_tmp.Unid)),
			}

			han_chars = append(han_chars, &han_char)

			//取最简字且常用字

			my_char := fate.Character{
				Ch: han_char.Ch,
			}

			has, err = GetFateChar(exc.Db, han_cheng_char_tmp.Unid, &my_char)

			if err != nil {
				continue
			}

			if has {
				my_chars = append(my_chars, &my_char)
			}
		}
	}

	for _, han_char := range han_chars {
		InsertOrUpdate(exc.Db, han_char)
	}

	for _, my_char := range my_chars {
		InsertOrUpdate(exc.DbFate, my_char)
	}

	return nil
}

//取最简字且常用字，必须有拼音、五行、姓名学笔画，不可以是偏旁字
func GetFateChar(exc_db *xorm.Engine, unid int, my_char *fate.Character) (bool, error) {

	pin_yin_id_gots := []models.PinYinId{}

	err := exc_db.Where("Unid = ?", unid).Find(&pin_yin_id_gots)
	if err != nil {
		panic(err)
	}

	pin_yin_gots := []string{}

	for _, pin_yin_id_got := range pin_yin_id_gots {
		pin_yin_got := models.PinYin{
			Pid: pin_yin_id_got.Pid,
		}

		exc_db.Get(&pin_yin_got)

		pin_yin_gots = append(pin_yin_gots, pin_yin_got.PinYin)
	}

	my_char.PinYin = pin_yin_gots

	if len(pin_yin_gots) > 1 {
		my_char.IsDuoYin = true
	} else if len(pin_yin_gots) == 1 {
		my_char.IsDuoYin = false
	} else {
		return false, errors.New("没有拼音")
	}

	min_su_id_got := models.MinSuId{
		Unid: unid,
	}

	has, err := exc_db.Get(&min_su_id_got)
	if err != nil {
		panic(err)
	}

	if !has {
		panic("民俗数据找不到")
	}

	min_su_got := models.MinSu{
		Msid: min_su_id_got.Msid,
	}

	exc_db.Get(&min_su_got)

	my_char.IsSurname = min_su_got.IsSurname
	my_char.SurnameGender = min_su_got.SurnameGender
	my_char.WuXing = min_su_got.WuXing
	my_char.Lucky = min_su_got.Lucky
	my_char.IsRegular = min_su_got.Regular

	science_stroke_got := models.ScienceStroke{
		Unid: unid,
	}

	has, err = exc_db.Get(&science_stroke_got)
	if err != nil {
		panic(err)
	}

	if !has {
		panic("姓名学笔画找不到")
	}

	my_char.ScienceStroke = science_stroke_got.ScienceStroke

	glyph_got := models.Glyph{
		Unid: unid,
	}

	has, err = exc_db.Get(&glyph_got)
	if err != nil {
		panic(err)
	}

	if !has {
		panic("字形信息找不到")
	}

	if glyph_got.Stroke == 0 {
		if glyph_got.SimplifiedTotalStroke == 0 {
			if glyph_got.TraditionalTotalStroke == 0 {
				panic("没有总笔画信息")
			} else {
				panic("这是个繁体字")
			}
		} else {
			my_char.Radical = glyph_got.SimplifiedRadical
			my_char.Stroke = glyph_got.SimplifiedTotalStroke
		}
	} else {
		my_char.Radical = glyph_got.Radical
		my_char.Stroke = glyph_got.Stroke
	}

	if !glyph_got.AsRadical {
		return true, nil
	} else {
		return false, errors.New("这是偏旁部首")
	}
}
