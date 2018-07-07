package excavator

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

//getRootList run get root
func getRootList(r *Root) {
	doc, err := parseDocument(r.URL)
	if err != nil {
		panic(err)
	}

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
				fmt.Println("log", stroke, ch, href, b)
				r.Add(&Radical{
					Strokes: strconv.Itoa(i),
					Name:    ch,
					URL:     href,
				})
			}
		})
	})
}

func getRedicalList(r *Root) {
	if r.HasNext() {

	}
}

func getCharacterList(r *Root) {

}

//ParseDocument get the url result body
func parseDocument(url string) (*goquery.Document, error) {
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
