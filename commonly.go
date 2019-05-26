package excavator

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)


type TopCallback func(url string, ch *RootRadicalCharacter)

//RootCYZ 常用字
func RootCYZ(host string, cb TopCallback) error {
	url := strings.Join([]string{host, "z/zb/cc1.htm"}, "/")
	html, e := parseDocument(url)
	if e != nil {
		return e
	}
	html.Find(".bs_index3").Each(func(i int, s1 *goquery.Selection) {
		s1.Find("li").Each(func(i int, s2 *goquery.Selection) {
			a := s2.Find("a").Text()
			link, _ := s2.Find("a").Attr("href")
			pinyin, _ := s2.Find("a").Attr("title")
			cc := RootRadicalCharacter{
				Character: a,
				Link:      link,
				Pinyin:    strings.Split(pinyin, ","),
			}
			log.Infof("%+v", cc)
			cb(host, &cc)
		})
	})
	return nil
}

//CommonlyBase
func CommonlyBase(url string, character *RootRadicalCharacter) {
	url = url + character.Link
	html, e := parseDocument(url)
	if e != nil {
		log.Errorf("%s error with:%s", e, url)
		return
	}
	bc := StandardCharacter{
		Character: character.Character,
		NeedFix:   true,
		Data:      make(map[string]string),
	}

	html.Find(".tab").Each(func(i int, s1 *goquery.Selection) {
		log.Info(s1.Html())
		k := strings.TrimSpace(s1.Text())
		v, b := s1.Find("a").Attr("href")
		if b {
			bc.Data[k] = v
		} else {
			bc.Data[k] = character.Link
		}
	})
	//return &bc
}
