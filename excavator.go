package excavator

import (
	"strings"
)

func Get(url string) map[string][]string {
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
		s = strings.Replace(s, " bgcolor='#F4F5F9'  align=center", "", -1)
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

func GetKangXi(url string) []CharDetail {
	var ret []CharDetail
	TransformOff()

	SetFix(func(s string) string {
		s = strings.Replace(s, `<td align=middle width="10%" bgcolor=#F0E4E1`, "<td", -1)
		return s
	})

	//http://tool.httpcn.com/KangXi/BuShou.html
	//goquery.NewDocument("http://tool.httpcn.com/KangXi/BuShou.html")
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
			//var cl []string
			list1 := GetFileterCharList(url + "/" + key)
			//TransformOn()
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

func UpdateKangXi(url string,f func(CharDetail))  {
	TransformOff()

	SetFix(func(s string) string {
		s = strings.Replace(s, `<td align=middle width="10%" bgcolor=#F0E4E1`, "<td", -1)
		return s
	})

	//http://tool.httpcn.com/KangXi/BuShou.html
	//goquery.NewDocument("http://tool.httpcn.com/KangXi/BuShou.html")
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
			//var cl []string
			list1 := GetFileterCharList(url + "/" + key)
			//TransformOn()
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