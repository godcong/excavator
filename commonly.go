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

func dummyLog(c *StandardCharacter, i int, selection *goquery.Selection) (err error) {
	log.With("index", i).Info(selection.Text())
	return nil
}

func mandarinDictionary(ex *StandardCharacter, i int, selection *goquery.Selection) (err error) {

	selection.Find("div[class=gycd]").Each(func(i int, selection *goquery.Selection) {
		var explan MandarinExplanation
		selection.Find("div[class=pz]").Each(func(i int, selection *goquery.Selection) {
			log.Info(selection.Text())
			selection.Find("rt").Each(func(i int, selection *goquery.Selection) {
				switch i {
				case 0:
					explan.Phonetic = selection.Text()
				case 1:
					explan.Pinyin = selection.Text()
				}
			})
		})
		selection.Find("div[class=gycd-item]").Each(func(i int, selection *goquery.Selection) {
			explan.PartOfSpeech = selection.Find("span[class=gc_cx]").Text()
			selection.Find("li").Each(func(i int, selection *goquery.Selection) {
				explan.Explanation = append(explan.Explanation, selection.Text())
			})
		})
		log.Infof("%+v", explan)
		ex.MandarinDictionary.Explanation = append(ex.MandarinDictionary.Explanation, explan)
	})

	log.Infof("%+v", ex.MandarinDictionary)
	return nil
}
func detailedExplanation(ex *StandardCharacter, i int, selection *goquery.Selection) (err error) {
	var explan DetailedExplanation
	selection.Find("div[class~=xnr]").Each(func(i int, selection *goquery.Selection) {

		selection.Find("hr[class=dichr]").Next().Each(func(i int, selection *goquery.Selection) {
			explan.Pinyin = selection.Find("span[class=dicpy]").Text()
			cino := selection.Find("span[class=cino]").Parent().Text()
			//TODO:FIX
			explan.DetailedMeaning = append(explan.DetailedMeaning, cino)
			dic := selection.Find("span[class=diczx1]").Parent().Text()
			explan.DetailedMeaning = append(explan.DetailedMeaning, dic)
		})
	})

	//ex.DetailedExplanation.DetailedMeaning = append(ex.DetailedExplanation.DetailedMeaning, selection.Text())
	//})
	//selection.Find()
	log.Infof("%+v", ex.DetailedExplanation)
	return nil
}

func basicExplanation(ex *StandardCharacter, i int, selection *goquery.Selection) (err error) {
	selection.Find("div[class~=jnr]").Each(func(i int, selection *goquery.Selection) {
		log.With("basicExplanation", "jnr").Debug(selection.Text())
		py := selection.Find("span[class=dicpy]").Text()
		s := strings.Split(py, " ")
		ex.BasicExplanation.Pinyin = strings.TrimSpace(s[0])
		ex.BasicExplanation.Phonetic = strings.TrimSpace(s[len(s)-1])
		log.With("basicExplanation", "pinyin").Debug(ex.BasicExplanation.Pinyin, ",", ex.BasicExplanation.Phonetic)
		selection.Find("ol li").Each(func(i int, selection *goquery.Selection) {
			log.With("index", i).Info(selection.Text())
			ex.BasicExplanation.BasicMeaning = append(ex.BasicExplanation.BasicMeaning, selection.Text())
		})
	})

	//selection.Find()
	log.Infof("%+v", ex.BasicExplanation)
	return nil
}

// ProcessFunc ...
type ProcessFunc func(*StandardCharacter, int, *goquery.Selection) error

// DataTypeBlock ...
var DataTypeBlock = map[string]ProcessFunc{
	"基本解释": basicExplanation,
	"详细解释": detailedExplanation,
	"國語詞典": mandarinDictionary,
	"康熙字典": dummyLog,
	"说文解字": dummyLog,
	"音韵方言": dummyLog,
	"字源字形": dummyLog,
	"网友讨论": dummyLog,
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
		log.Debug(s1.Html())

		dtb, b := s1.Attr("data-type-block")
		if !b {
			return
		}
		if fn, b := DataTypeBlock[dtb]; b {
			e := fn(&bc, i, s1)
			if e != nil {
				log.Error(e)
				return
			}

		}

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
