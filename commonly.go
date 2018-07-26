package excavator

import (
		"github.com/PuerkitoBio/goquery"
		)

//CommonlyTop
func CommonlyTop(url string) []string {
	var chars []string
	html, _ := parseDocument(url)
	html.Find(".bs_index3").Each(func(i int, s1 *goquery.Selection) {
		s1.Find("li").Each(func(i int, s2 *goquery.Selection) {
			val := s2.Find("a").Text()
			chars = append(chars, val)
		})
	})
	return chars
}
