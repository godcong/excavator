package excavator

import (
	"log"
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
	log.Println("map:", retMap)
	return retMap
}
