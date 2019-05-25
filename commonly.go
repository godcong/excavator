package excavator

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/globalsign/mgo/bson"
)

type RootRadicalCharacter struct {
	Character string   `json:"character"`
	Link      string   `json:"link"`
	Pinyin    []string `json:"pinyin"`
}

type BaseCharacter struct {
	ID        bson.ObjectId `bson:"_id,omitempty"`
	NeedFix   bool          `bson:"need_fix"`
	Character string
	Data      map[string]string
}

type TopCallback func(url string, ch *RootRadicalCharacter)

//RootCYZ
func RootCYZ(url string, cb TopCallback) error {
	url = strings.Join([]string{url, "z/zb/cc1.htm"}, "/")
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
			cb(url, &cc)
		})
	})
	return nil
}

//CommonlyBase
func CommonlyBase(url string, character *RootRadicalCharacter) {
	url = url + character.Link
	html, _ := parseDocument(url)
	//bc := StandardCharacter{
	//	Character: character.Character,
	//	NeedFix:   true,
	//	Data:      make(map[string]string),
	//}
	//if err != nil {
	//	return &bc
	//}

	html.Find(".tab").Each(func(i int, s1 *goquery.Selection) {
		fmt.Println(s1.Html())
		//k := strings.TrimSpace(s1.Text())
		//v, b := s1.Find("a").Attr("href")
		//if b {
		//	bc.Data[k] = v
		//} else {
		//	bc.Data[k] = character.Link
		//}
	})
	//return &bc
}
