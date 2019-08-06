package excavator

import (
	"github.com/go-xorm/xorm"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

// CharacterFunc ...
type CharacterFunc func(character *Character) error

//Character 字符
type Character struct {
	PinYin                   []string `xorm:"default() notnull pin_yin"`                     //拼音
	Character                string   `xorm:"default() notnull character"`                   //字符
	Radical                  string   `xorm:"default() notnull radical"`                     //部首
	RadicalStroke            int      `xorm:"default(0) notnull radical_stroke"`             //部首笔画
	IsKangXi                 bool     `xorm:"default(0) notnull is_kang_xi"`                 //是否康熙字典
	KangXi                   string   `xorm:"default() notnull kang_xi"`                     //康熙
	KangXiStroke             int      `xorm:"default(0) notnull kang_xi_stroke"`             //康熙笔画
	SimpleRadical            string   `xorm:"default() notnull simple_radical"`              //简体部首
	SimpleRadicalStroke      int      `xorm:"default(0) notnull simple_radical_stroke"`      //简体部首笔画
	SimpleTotalStroke        int      `xorm:"default(0) notnull simple_total_stroke"`        //简体笔画
	TraditionalRadical       string   `xorm:"default() notnull traditional_radical"`         //繁体部首
	TraditionalRadicalStroke int      `xorm:"default(0) notnull traditional_radical_stroke"` //繁体部首笔画
	TraditionalTotalStroke   int      `xorm:"default(0) notnull traditional_total_stroke"`   //简体部首笔画
	NameScience              bool     `xorm:"default(0) notnull name_science"`               //姓名学
	WuXing                   string   `xorm:"default() notnull wu_xing"`                     //五行
	Lucky                    string   `xorm:"default() notnull lucky"`                       //吉凶寓意
	Regular                  bool     `xorm:"default(0) notnull regular"`                    //常用
	TraditionalCharacter     []string `xorm:"default() notnull traditional_character"`       //繁体字
	VariantCharacter         []string `xorm:"default() notnull variant_character"`           //异体字
	Comment                  []string `xorm:"default() notnull comment"`                     //解释
}

// ParseFunc ...
type ParseFunc func(*Character, int, string)

var charList = map[string]ParseFunc{
	"部首:":     parseBuShou,
	"简体部首:":   parseSimple,
	"繁体部首:":   parseTraditional,
	"康熙字典笔画:": parseKangXi,
	"拼音":      parsePinYin,
}

// Clone ...
func (c *Character) Clone() (char *Character) {
	char = new(Character)
	*char = *c
	copy(char.PinYin, c.PinYin)
	return char
}

// InsertIfNotExist ...
func (c *Character) InsertIfNotExist(session *xorm.Session) (e error) {

	i, e := session.Table(&Character{}).Where("character = ?", c.Character).Count()
	if e != nil {
		return e
	}
	if i == 0 {
		_, e = session.InsertOne(c)
		return
	}
	_, e = session.Where("character = ?", c.Character).Update(c)
	//if e != nil {
	//	return e
	//}
	return
}

func parseDummy(c *Character, index int, input string) {
	log.With("character", c, "index", index, "input", input).Info("dummy")
}
func parseKangXi(c *Character, index int, input string) {
	log.With("character", c, "index", index, "input", input).Info("kangxi")
	log.With("input", input).Info("kangxi bracket")
	input = strings.ReplaceAll(input, "(", " ")
	input = strings.ReplaceAll(input, ")", " ")
	s := strings.Split(strings.TrimSpace(input), ";")
	if c.KangXiStroke == 0 && len(s) > 0 {
		vv := strings.Split(strings.TrimSpace(s[0]), ":")
		if len(vv) == 2 {
			c.KangXi = vv[0]
			i, e := strconv.Atoi(strings.TrimSpace(vv[1]))
			if e != nil {
				log.Error(e)
				return
			}
			c.KangXiStroke = i
		}
	}
}
func parseBuShou(c *Character, index int, input string) {
	log.With("input", input).Info("bushou")
	switch index {
	case 0:
		c.Radical = input
	case 1:
		parseNumber(&c.RadicalStroke, input)
	case 2:
		parseNumber(&c.KangXiStroke, input)
	default:
		log.Error("bushou")
	}
}

func parseSimple(c *Character, index int, input string) {
	log.With("input", input).Info("simple")
	switch index {
	case 0:
		parseBuShouBracket(input, &c.SimpleRadical, &c.SimpleRadicalStroke)
	case 1:
		parseNumber(&c.SimpleTotalStroke, input)
	default:
		log.Error("simple")
	}
}

func parseTraditional(c *Character, index int, input string) {
	log.With("input", input).Info("traditional")
	switch index {
	case 0:
		parseBuShouBracket(input, &c.TraditionalRadical, &c.TraditionalRadicalStroke)
	case 1:
		parseNumber(&c.TraditionalTotalStroke, input)
	default:
		log.Error("traditional")
	}
}

func newDoc(element *colly.HTMLElement) (d *goquery.Document, e error) {
	html, e := element.DOM.Html()
	if e != nil {
		log.Error(e)
		return
	}
	d, e = goquery.NewDocumentFromReader(strings.NewReader(html))
	return
}

func parseKangXiCharacter(element *colly.HTMLElement, ch *Character) (e error) {
	n, e := newDoc(element)
	if e != nil {
		log.Error(e)
		return e
	}
	v := StringClearUp(n.ReplaceWith("font[class=colred]").Text())
	var data []string
	if v != "" {
		data = strings.Split(v, " ")
	}

	log.With("source", v).Info(len(data), ":", data)
	n1, e := newDoc(element)
	if e != nil {
		log.Error(e)
		return e
	}
	f := parseDummy
	n1.Find("font[class=colred]").Each(func(i int, selection *goquery.Selection) {
		log.With("text", selection.Text(), "index", selection.Index(), "num", i).Info("colred")
		text := StringClearUp(selection.Text())
		if i == 0 {
			if data == nil || len(data) == 0 {
				f = parsePinYin
			} else {
				if v, b := charList[text]; b {
					f = v
				}
			}
		}
		if len(data) > i {
			f(ch, i, data[i])
		} else {
			if data == nil || len(data) == 0 {
				f(ch, i, text)
			}
		}
	})
	return nil
}

func parseArray(source *[]string, input string) {
	*source = append(*source, input)
}
func parseNumber(source *int, input string) {
	i, e := strconv.Atoi(strings.ReplaceAll(input, "画", ""))
	if e != nil {
		log.With("input", input).Error(e)
		return
	}
	*source = i
}
func parsePinYin(c *Character, index int, input string) {
	log.With("input", input).Info("pinyin")
	input = strings.ReplaceAll(input, "[", "")
	input = strings.ReplaceAll(input, "]", "")
	parseArray(&c.PinYin, input)
}
func parseBuShouBracket(input string, radical *string, stroke *int) {
	log.With("input", input).Info("bushou bracket")
	input = strings.ReplaceAll(input, "(", " ")
	input = strings.ReplaceAll(input, ")", " ")
	input = strings.TrimSpace(input)
	s := strings.Split(input, " ")
	if len(s) == 2 {
		*radical = s[0]
		*stroke, _ = strconv.Atoi(s[1])
	}
}

var infoList = map[string]ParseFunc{
	"汉字五行：":   parseWuXing,
	"吉凶寓意：":   parseLucy,
	"姓名学：":    parseNameScience,
	"是否为常用字：": parseRegular,
	//"繁体字集：":   parseTraditionalCharacter,
	//"异体字集：":   parseVariantCharacter,
}
var infoList2 = map[string]ParseFunc{
	"繁体字集：": parseTraditionalCharacter,
	"异体字集：": parseVariantCharacter,
}

func parseVariantCharacter(c *Character, index int, input string) {
	log.With("input", input).Info("var char")
	if input != "" {
		c.VariantCharacter = append(c.VariantCharacter, input)
	}
}
func parseTraditionalCharacter(c *Character, index int, input string) {
	log.With("input", input).Info("trad char")
	if input != "" {
		c.TraditionalCharacter = append(c.TraditionalCharacter, input)
	}
}
func parseRegular(c *Character, index int, input string) {
	log.With("input", input).Info("regular")
	if input == "是" {
		c.Regular = true
	}
}
func parseNameScience(c *Character, index int, input string) {
	log.With("input", input).Info("name science")
	if input == "是" {
		c.NameScience = true
	}
}
func parseLucy(c *Character, index int, input string) {
	log.With("input", input).Info("lucky")
	c.Lucky = input
}
func parseWuXing(c *Character, index int, input string) {
	log.With("input", input).Info("wuxing")
	c.WuXing = input
}
func parseDictInformation(element *colly.HTMLElement, ch *Character) (e error) {
	fn := parseDummy
	n, e := newDoc(element)
	if e != nil {
		log.Error(e)
		return e
	}

	n.Find("li").Contents().Each(func(i int, selection *goquery.Selection) {
		log.With("text", selection.Text(), "index", selection.Index(), "num", i).Info("li")
		tx := selection.Text()
		if selection.Index() == 0 {
			fn = parseDummy
			if v, b := infoList[tx]; b {
				fn = v
			}
			return
		}

		if goquery.NodeName(selection) == "#text" {
			fn(ch, i, StringClearUp(tx))
		}
	})

	n1, e := newDoc(element)
	if e != nil {
		log.Error(e)
		return e
	}
	n1.Find("li").Each(func(i int, selection *goquery.Selection) {
		log.With("text", selection.Text(), "index", selection.Index(), "num", i).Info("li2")
		tx := selection.Find("span").Text()
		if v, b := infoList2[tx]; b {
			v(ch, selection.Index(), selection.Find("a").Text())
		}
	})
	return
}
func parseComment(c *Character, index int, input string) {
	log.With("input", input).Info("comment")
	if input == "" {
		return
	}
	c.Comment = append(c.Comment, StringClearUp(input))
}

func parseDictComment(element *colly.HTMLElement, character *Character) (e error) {
	n, e := newDoc(element)
	if e != nil {
		log.Error(e)
		return e
	}
	tx := StringClearUp(n.Find("li > a").Text())
	if tx == "康熙字典解释" {
		n.Find("li > div").Contents().Each(func(i int, selection *goquery.Selection) {
			log.With("text", selection.Text(), "index", selection.Index(), "num", i).Info("li3")
			parseComment(character, i, selection.Text())
		})
	}

	return
}
