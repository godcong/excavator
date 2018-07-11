package excavator

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type CharacterAssign func(c *Character, s string) bool

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
							TotalStrokes:   "",
							KangxiStrokes:  ch[2],
							Phonetic:       "",
							Folk:           Folk{},
							Structure:      Structure{},
						})
					}
				}
			})

		})
	})
}

func saveRadical(c *Character, v string) bool {
	log.Println(v)
	c.Radical = v
	return true
}

func saveRadicalStrokes(c *Character, v string) bool {
	log.Println(v)
	c.RadicalStrokes = v
	return true
}

func saveTotalStrokes(c *Character, v string) bool {
	log.Println(v)
	c.TotalStrokes = v
	return true
}

func saveKangXiStrokes(c *Character, v string) bool {
	log.Println(v)
	c.KangxiStrokes = v
	return true
}

func savePinyin(c *Character, v string) bool {
	log.Println(v)
	c.Pinyin = v
	return true
}

func savePhonetic(c *Character, v string) bool {
	log.Println(v)
	c.Phonetic = v
	return true
}

func dummySave(c *Character, v string) bool {
	log.Println(v)
	return true
}

var caMap = map[string]CharacterAssign{
	"部首笔画：":  saveRadicalStrokes,
	"部首：":    saveRadical,
	"总笔画：":   saveTotalStrokes,
	"康熙字典笔画": saveKangXiStrokes,
	"拼音：":    savePinyin,
	"注音：":    savePhonetic,
}

func characterAssignFunc(list map[string]CharacterAssign, text string) CharacterAssign {
	for k, v := range list {
		if strings.Index(text, k) >= 0 {
			return v
		}
	}
	return nil
}

func folk(c *Character, text string) bool {
	s := strings.Split(text, "：")
	if strings.Index(text, "是否为常用字") >= 0 {
		if len(s) >= 2 {
			c.Folk.CommonlyCharacters = s[1]
		}

		return true
	} else if strings.Index(text, "姓名学") >= 0 {
		if len(s) >= 2 {
			c.Folk.NameScience = s[1]
		}
		return true
	}
	return false
}

func decompositionSearch(c *Character, v string) bool {
	c.DecompositionSearch = strings.Replace(v, "]：", "", -1)
	return true
}

func strokeNumber(c *Character, v string) bool {
	c.StrokeNumber = strings.Replace(v, "]：", "", -1)
	return true
}
func strokeReadWrite(c *Character, v string) bool {
	c.StrokeReadWrite = strings.Replace(v, "]：", "", -1)
	return true
}

var sMap = map[string]CharacterAssign{
	"首尾分解查字": decompositionSearch,
	"笔顺编号":   strokeNumber,
	"笔顺读写":   strokeReadWrite,
	"广　韵":    dummySave,
	"国　语":    dummySave,
	"粤　语":    dummySave,
	"古文字诂林":  dummySave,
	"故训彙纂":   dummySave,
	"说文解字":   dummySave,
	"康熙字典":   dummySave,
	"汉语字典":   dummySave,
	"辞　海":    dummySave,
}

func getCharacterList(r *Root, c *Character) {
	url := r.URL + c.URL
	doc, err := parseDocument(url)
	if err != nil {
		panic(err)
	}
	log.Println(doc.Html())
	docCopy := doc.Clone()

	//处理笔画
	f := dummySave
	doc.Find("tbody tr .text15").ReplaceWith("script").Each(func(i int, selection *goquery.Selection) {
		selection.Contents().Each(func(i int, selection *goquery.Selection) {
			text := selection.Text()

			if f != nil {
				if b := f(c, text); b {
					f = nil
				}
			}

			f = characterAssignFunc(caMap, text)

		})
	})

	docCopy.Find(".text16").Each(func(i int, selection *goquery.Selection) {
		selection.Contents().Each(func(i int, selection *goquery.Selection) {
			text := selection.Text()

			if b := folk(c, text); b {
				return
			}

			if f != nil {
				if b := f(c, text); b {
					f = nil
				}
			}

			f = characterAssignFunc(sMap, text)
		})

	})
	log.Println(*c)

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
