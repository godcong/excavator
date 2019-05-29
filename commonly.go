package excavator

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// TopCallback ...
type TopCallback func(url string, ch *RootRadicalCharacter)

//RootRegular 常用字
func RootRegular(host string, cb TopCallback) error {
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
				Class:     ClassRegular,
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

// DataTypeBlock ...
var DataTypeBlock = []string{
	"基本解释",
	"详细解释",
	"國語詞典",
	"康熙字典",
	"说文解字",
	"音韵方言",
	"字源字形",
	"网友讨论",
}

//CommonlyBase ...
func CommonlyBase(url string, character *RootRadicalCharacter) {
	url = url + character.Link
	html, e := parseDocument(url)
	if e != nil {
		log.Errorf("%s error with:%s", e, url)
		return
	}
	bc := StandardCharacter{
		Radical:         character.Character,
		CharacterDetail: map[string]string{},
	}

	html.Find("div[data-type-block]").Each(func(i int, s1 *goquery.Selection) {
		//DataTypeBlock
		log.Info(s1.Attr("data-type-block"))
		//log.Info(s1.Text())
		//k := strings.TrimSpace(s1.Text())
		//v, b := s1.Find("data-type-block").Attr("href")
		//if !b || k == "基本解释" {
		//	//基本解释
		//	bc.CharacterDetail[k] = character.Link
		//	return
		//}
		////other
		//bc.CharacterDetail[k] = v
	})
	log.Infof("character detail:%+v", bc)
}
