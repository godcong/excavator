package excavator

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type CharacterAssign func(c *Character, s string) bool

//getRootList run get root
func getRootList(r *Root, suffix string) *Root {
	doc, err := parseDocument(r.URL + suffix)
	if err != nil {
		log.Info(r.URL + suffix)
		panic(err)
	}

	doc.Find("table tbody").Each(func(i int, s *goquery.Selection) {
		stroke := i
		s.Find("tr td").Each(func(i1 int, selection *goquery.Selection) {
			if i == 0 {
				return
			}

			href, b := selection.Find("a").Attr("href")
			ch := selection.Text()
			if b {
				rc := &RootCharacter{
					Strokes: strconv.Itoa(stroke),
					Name:    trim(ch),
					URL:     href,
				}
				r.Add(rc)
			}
		})
	})
	return r
}

func getRedicalList(r *Root, rc *RootCharacter) *Radical {
	radical := new(Radical)
	url := r.URL + rc.URL
	doc, err := parseDocument(url)
	if err != nil {
		log.Info(url)
		return radical
	}

	doc.Find("table tbody").Each(func(i int, selection *goquery.Selection) {
		selection.Find("tr").Each(func(i1 int, selectiontr *goquery.Selection) {
			if i1 == 0 {
				return
			}
			ch := make([]string, 5)
			selectiontr.Find("td").Each(func(i2 int, selection *goquery.Selection) {
				html, _ := selection.Html()
				html = trim(html)
				add := false
				switch i2 % 4 {
				case 1:
					ch[0] = html
				case 2:
					ch[1] = html
				case 3:
					ch[2] = html
					add = true
				case 0:
					ch[3] = selection.Find("a").Text()
					href, b := selection.Find("a").Attr("href")
					if b {
						ch[4] = href
					}
				}

				if add {
					if ch[3] == "" {
						log.Info(selectiontr.Html())
						return
					}

					radicalCharacter := &RadicalCharacter{
						RootCharacter: rc,
						Strokes:       ch[2],
						Pinyin:        ch[0],
						Character:     ch[3],
						URL:           ch[4],
					}
					radical.Add(radicalCharacter)
					ch = make([]string, 5)
				}
			})

		})
	})
	return radical
}

func saveRadical(c *Character, v string) bool {
	return trimReplace(&c.Radical, v)
}

func saveRadicalStrokes(c *Character, v string) bool {
	return trimReplace(&c.RadicalStrokes, v)
}

func saveTotalStrokes(c *Character, v string) bool {
	return trimReplace(&c.TotalStrokes, v)
}

func saveKangXiStrokes(c *Character, v string) bool {
	return trimReplace(&c.KangxiStrokes, v)
}

func savePinyin(c *Character, v string) bool {
	v = strings.Replace(v, "　", ",", -1)
	return trimReplace(&c.Pinyin, v)
}

func savePhonetic(c *Character, v string) bool {
	return trimReplace(&c.Phonetic, v)
}

func dummySave(c *Character, v string) bool {
	log.Info(v)
	return true
}

func characterSave(c *Character, v string) bool {
	v = strings.Replace(v, "『", "", -1)
	v = strings.Replace(v, "』", "", -1)
	return trimReplace(&c.Character, v)
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

func indexFilter(text string, subs []string) []int {
	var idxs []int
	for _, v := range subs {
		if i := strings.Index(text, v); i >= 0 {
			idxs = append(idxs, i)
		}
	}
	sort.Ints(idxs)
	return idxs
}

func folk(c *Character, text string) bool {
	subs := []string{"是否为常用字", "姓名学", "汉字五行", "吉凶寓意"}
	var rlts []string
	idxs := indexFilter(text, subs)

	switch len(idxs) {
	case 4:
		rlts = append(rlts, strings.TrimSpace(text[idxs[0]:idxs[1]]))
		rlts = append(rlts, strings.TrimSpace(text[idxs[1]:idxs[2]]))
		rlts = append(rlts, strings.TrimSpace(text[idxs[2]:idxs[3]]))
		rlts = append(rlts, strings.TrimSpace(text[idxs[3]:]))
	case 3:
		rlts = append(rlts, strings.TrimSpace(text[idxs[0]:idxs[1]]))
		rlts = append(rlts, strings.TrimSpace(text[idxs[1]:idxs[2]]))
		rlts = append(rlts, strings.TrimSpace(text[idxs[2]:]))
	case 2:
		rlts = append(rlts, strings.TrimSpace(text[idxs[0]:idxs[1]]))
		rlts = append(rlts, strings.TrimSpace(text[idxs[1]:]))
	case 1:
		rlts = append(rlts, strings.TrimSpace(text[idxs[0]:]))
	case 0:
		return false
	}
	for _, v := range rlts {
		v1 := strings.Split(v, "：")
		if len(v1) >= 2 {

			switch v1[0] {
			case "是否为常用字":
				c.CommonlyCharacters = trim(v1[1])
			case "姓名学":
				c.NameScience = trim(v1[1])
			case "汉字五行":
				c.FiveElementCharacter = trim(v1[1])
			case "吉凶寓意":
				c.GodBadMoral = trim(v1[1])
			}
		}
	}

	return true
}

func trimReplace(source *string, s string) bool {
	if s == "" {
		return false
	}
	*source = trim(strings.Replace(s, "]：", "", -1))
	return true
}

func decompositionSearch(c *Character, v string) bool {
	return trimReplace(&c.DecompositionSearch, v)
}

func strokeNumber(c *Character, v string) bool {
	return trimReplace(&c.StrokeNumber, v)
}
func strokeReadWrite(c *Character, v string) bool {
	return trimReplace(&c.StrokeReadWrite, v)
}

//guangYun 广　韵
func guangYun(c *Character, v string) bool {
	return trimReplace(&c.GuangYun, v)
}

//国　语
func mandarin(c *Character, v string) bool {
	return trimReplace(&c.Mandarin, v)
}

////粤　语
func cantonese(c *Character, v string) bool {
	return trimReplace(&c.Cantonese, v)
}

//古文字诂林
func ancientWrite(c *Character, v string) bool {
	return trimReplace(&c.AncientWrite, v)
}

//故训彙纂
func hometownTrain(c *Character, v string) bool {
	return trimReplace(&c.HometownTrain, v)
}

//说文解字
func explain(c *Character, v string) bool {
	return trimReplace(&c.Index.Explain, v)
}

//康熙字典
func kangxiDictionary(c *Character, v string) bool {
	return trimReplace(&c.KangxiDictionary, v)
}

//汉语字典
func chineseDictionary(c *Character, v string) bool {
	return trimReplace(&c.ChineseDictionary, v)
}

//辞　海  
func cihai(c *Character, v string) bool {
	return trimReplace(&c.Cihai, v)
}

var sMap = map[string]CharacterAssign{
	"首尾分解查字": decompositionSearch,
	"笔顺编号":   strokeNumber,
	"笔顺读写":   strokeReadWrite,
	"广　韵":    guangYun,
	"国　语":    mandarin,
	"粤　语":    cantonese,
	"古文字诂林":  ancientWrite,
	"故训彙纂":   hometownTrain,
	"说文解字":   explain,
	"康熙字典":   kangxiDictionary,
	"汉语字典":   chineseDictionary,
	"辞　海":    cihai,
}

func getCharacterList(r *Root, rc *RadicalCharacter) *Character {
	c := new(Character)
	url := r.URL + rc.URL
	doc, err := parseDocument(url)
	if err != nil {
		log.Info(err, url)
		return &Character{
			Character:      rc.Character,
			Pinyin:         rc.Pinyin,
			Radical:        "",
			RadicalStrokes: "",
			TotalStrokes:   rc.Strokes,
			KangxiStrokes:  "",
			Phonetic:       "",
			Folk:           Folk{},
			Structure:      Structure{},
			Explain:        Explain{},
			Rhyme:          Rhyme{},
			Index:          Index{},
		}
	}

	//处理笔画
	f := characterSave
	doc.Find("tbody tr .text15").ReplaceWith("script").Each(func(i int, s1 *goquery.Selection) {

		s1.Contents().Each(func(i int, selection *goquery.Selection) {
			text := selection.Text()
			if f != nil {
				if b := f(c, text); b {
					f = nil
				}
				//log.Println(text)
				//log.Println(s1.Html())
			}
			f = characterAssignFunc(caMap, text)
		})
	})
	kxExplain := doc.Clone()
	kxExplain.Find(".content16").Each(func(i int, selection *goquery.Selection) {
		c.Explain.Intro, _ = selection.Find("strong").Html()
		c.Explain.Detail = selection.Text()
	})

	docCopy := doc.Clone()
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
	return c
}

//ParseDocument get the url result body
func parseDocument(url string) (doc *goquery.Document, e error) {
	var reader io.Reader
	hash := MD5(url)
	if !CheckExist(hash) {
		// Request the HTML page.
		res, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()
		if res.StatusCode != 200 {
			return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
		}
		reader = res.Body
		file, e := os.OpenFile(GetPath(hash), os.O_CREATE|os.O_SYNC|os.O_RDWR, os.ModePerm)
		if e != nil {
			return nil, e

		}
		written, e := io.Copy(file, reader)
		if e != nil {
			return nil, e
		}
		log.Infof("read %s | %d ", hash, written)
		_ = file.Close()
	}
	reader, e = os.Open(GetPath(hash))
	if e != nil {
		return nil, e
	}
	// Load the HTML document
	return goquery.NewDocumentFromReader(reader)
}

func trim(s string) string {
	s = strings.TrimSpace(s)
	s = strings.Replace(s, "　", "", -1)
	s = strings.Replace(s, "\n", "", -1)
	return s
}
