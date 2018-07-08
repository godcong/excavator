package excavator

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

//getRootList run get root
func getRootList(r *Root, suffix string) {
	doc, err := parseDocument(r.URL + suffix)
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
				r.Add(&Radical{
					Strokes: strconv.Itoa(stroke),
					Name:    ch,
					URL:     href,
				})
			}
		})
	})
}

func getRedicalList(r *Root, radical *Radical) {
	url := r.URL + radical.URL
	doc, err := parseDocument(url)
	if err != nil {
		panic(err)
	}
	doc.Find("table tbody").Each(func(i int, selection *goquery.Selection) {
		selection.Find("tr").Each(func(i1 int, selection *goquery.Selection) {

			if i1 == 0 {
				return
			}
			ch := make([]string, 5)
			selection.Find("td").Each(func(i2 int, selection *goquery.Selection) {
				html, _ := selection.Html()
				switch i2 % 4 {
				case 1:
					ch[0] = html
				case 2:
					ch[1] = html
				case 3:
					ch[2] = html
				case 0:
					ch[3] = selection.Find("a").Text()
					href, b := selection.Find("a").Attr("href")
					if b {
						ch[4] = href
						radical.Add(&Character{
							URL:            ch[4],
							Character:      ch[3],
							Pinyin:         ch[0],
							Radical:        ch[1],
							RadicalStrokes: radical.Strokes,
							KangxiStrokes:  ch[2],
							Phonetic:       "",
							Folk:           Folk{},
						})
					}
				}
			})

		})
	})
}

func getCharacterList(r *Root, c *Character) {
	url := r.URL + c.URL
	doc, err := parseDocument(url)
	if err != nil {
		panic(err)
	}
	log.Println(doc.Html())

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
