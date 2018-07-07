package excavator

import (
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

//GetRootList run get root
func GetRootList(url string) {
	doc, err := ParseDocument(url)
	if err != nil {
		panic(err)
	}
	//root := NewRoot()

	doc.Find("table tbody").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		stroke := i
		//s.Find("td").Each(func(i int, selection *goquery.Selection) {
		s.Find("tr td").Each(func(i1 int, selection *goquery.Selection) {
			if i == 0 {
				return
			}

			href, b := selection.Find("a").Attr("href")
			ch := selection.Text()
			if b {
				fmt.Println("log", stroke,ch, href, b)
				//root.Radicals
			}

		})

		//})

	})
}

//ParseDocument get the url result body
func ParseDocument(url string) (*goquery.Document, error) {
	// Request the HTML page.
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	body, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
