package excavator

import (
	"log"

	"github.com/PuerkitoBio/goquery"
)

type CommonlyCharacter struct {
	ID        bson.ObjectId `bson:"_id,omitempty"`
	Character string        `bson:"character"`
	Link      string        `bson:"link"`
	Title     string        `bson:"title"`
}

//CommonlyTop
func CommonlyTop(url string) []*CommonlyCharacter {
	var chars []*CommonlyCharacter
	html, _ := parseDocument(url)
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
			log.Printf("%+v", cc)
			chars = append(chars, &cc)

		})
	})
	return chars
}
