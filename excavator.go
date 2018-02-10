package excavator

import (
	"strconv"
	"strings"
)

func GetDictionary(url string) map[string][]string {
	retMap := make(map[string][]string)

	SetFix(func(s string) string {
		s = strings.Replace(s, "class=font_14", "", -1)
		s = strings.Replace(s, `""`, `"`, -1)
		return s
	})
	TransformOn()

	list := GetRootList(url + "/bs.html")

	SetFix(func(s string) string {
		s = strings.Replace(s, " bgcolor=#ffffff ", "", -1)
		s = strings.Replace(s, " class=font_14", "", -1)
		s = strings.Replace(s, " bgcolor='#f4f5f9' ", "", -1)
		s = strings.Replace(s, " align=center", "", -1)
		return s
	})

	for _, v := range list {
		for key, value := range v {
			var cl []string
			list1 := GetCharList(url + "/" + key)
			for _, value1 := range list1 {
				for _, value2 := range value1 {
					cl = append(cl, value2)
				}
			}
			retMap[value] = cl

		}

	}
	return retMap
}

func UpdateDictionary(url string, f func(detail CharDetail)) {
	SetFix(func(s string) string {
		s = strings.Replace(s, "class=font_14", "", -1)
		s = strings.Replace(s, `""`, `"`, -1)
		return s
	})
	TransformOn()

	list := GetRootList(url + "/bs.html")

	SetFix(func(s string) string {
		s = strings.Replace(s, " bgcolor=#ffffff ", "", -1)
		s = strings.Replace(s, " class=font_14", "", -1)
		s = strings.Replace(s, " bgcolor='#f4f5f9' ", "", -1)
		s = strings.Replace(s, " align=center", "", -1)
		return s
	})

	for _, v := range list {
		for key, value := range v {
			list1 := GetCharList(url + "/" + key)
			for _, value1 := range list1 {
				for _, value2 := range value1 {
					//cl = append(cl, value2)
					s := strings.Split(value2, "|")
					if len(s) > 2 {
						strokes, _ := strconv.Atoi(s[0])
						c := CharDetail{
							Char:           s[1],
							NameType:       "",
							NameRoot:       "",
							Pinyin:         s[2],
							Radical:        value,
							SimpleStrokes:  strokes,
							ScienceStrokes: 0,
						}
						f(c)
					}
				}
			}
			//retMap[value] = cl

		}

	}
}

func GetKangXi(url string) []CharDetail {
	var ret []CharDetail
	TransformOff()

	SetFix(func(s string) string {
		s = strings.Replace(s, `<td align=middle width="10%" bgcolor=#F0E4E1`, "<td", -1)
		return s
	})

	list := GetRootList(url + "/KangXi/BuShou.html")

	SetFix(func(s string) string {
		s = strings.Replace(s, "<tr bgcolor=#ffffff>", "<tr>", -1)

		s = strings.Replace(s, `<td align="center" bgcolor="#f7f1f0">`, "<td>", -1)
		s = strings.Replace(s, `<font color=red size=4>`, "<font>", -1)
		s = strings.Replace(s, `<a title=点击显示注释 href=`, `<a href="`, -1)
		s = strings.Replace(s, `><font`, `"><font`, -1)
		return s
	})

	for _, v := range list {
		for key := range v {
			list1 := GetFileterCharList(url + "/" + key)
			for _, v1 := range list1 {
				for k := range v1 {
					//cl = append(cl, value2)
					d := GetCharDetail(url + k)
					if d.Char != "" {
						ret = append(ret, d)
					}
				}
			}
		}
	}
	return ret
}

func UpdateKangXi(url string, f func(CharDetail)) {
	TransformOff()

	SetFix(func(s string) string {
		s = strings.Replace(s, `<td align=middle width="10%" bgcolor=#F0E4E1`, "<td", -1)
		return s
	})

	list := GetRootList(url + "/KangXi/BuShou.html")

	SetFix(func(s string) string {
		s = strings.Replace(s, "<tr bgcolor=#ffffff>", "<tr>", -1)

		s = strings.Replace(s, `<td align="center" bgcolor="#f7f1f0">`, "<td>", -1)
		s = strings.Replace(s, `<font color=red size=4>`, "<font>", -1)
		s = strings.Replace(s, `<a title=点击显示注释 href=`, `<a href="`, -1)
		s = strings.Replace(s, `><font`, `"><font`, -1)
		return s
	})

	for _, v := range list {
		for key := range v {
			list1 := GetFileterCharList(url + "/" + key)
			for _, v1 := range list1 {
				for k := range v1 {
					//cl = append(cl, value2)
					d := GetCharDetail(url + k)
					if d.Char != "" {
						f(d)
					}
				}
			}
		}
	}
}
