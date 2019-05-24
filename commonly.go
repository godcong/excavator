package excavator

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/globalsign/mgo/bson"
)

type CommonlyCharacter struct {
	ID        bson.ObjectId `bson:"_id,omitempty"`
	Character string        `bson:"character"`
	Link      string        `bson:"link"`
	Title     string        `bson:"title"`
}

type BaseCharacter struct {
	ID        bson.ObjectId `bson:"_id,omitempty"`
	NeedFix   bool          `bson:"need_fix"`
	Character string
	Data      map[string]string
}

//CommonlyTop
func CommonlyTop(url string) []*CommonlyCharacter {
	var chars []*CommonlyCharacter
	html, e := parseDocument(url)
	if e != nil {
		return nil
	}
	html.Find(".bs_index3").Each(func(i int, s1 *goquery.Selection) {
		s1.Find("li").Each(func(i int, s2 *goquery.Selection) {
			a := s2.Find("a").Text()
			link, _ := s2.Find("a").Attr("href")
			title, _ := s2.Find("a").Attr("title")
			cc := CommonlyCharacter{
				Character: a,
				Link:      link,
				Title:     title,
			}
			log.Infof("%+v", cc)
			chars = append(chars, &cc)

		})
	})
	return chars
}

//CommonlyBase
func CommonlyBase(url string, character *CommonlyCharacter) *BaseCharacter {
	url = url + character.Link
	html, err := parseDocument(url)
	bc := BaseCharacter{
		Character: character.Character,
		NeedFix:   true,
		Data:      make(map[string]string),
	}
	if err != nil {
		return &bc
	}

	html.Find(".tab").Each(func(i int, s1 *goquery.Selection) {
		fmt.Println(s1.Html())
		k := strings.TrimSpace(s1.Text())
		v, b := s1.Find("a").Attr("href")
		if b {
			bc.Data[k] = v
		} else {
			bc.Data[k] = character.Link
		}
	})
	return &bc
}
