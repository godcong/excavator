package excavator

import (
	"github.com/go-xorm/xorm"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// CharacterFunc ...
type CharacterFunc func(character *Character) error

//Character 字符
type Character struct {
	PinYin                   []string `xorm:"default() notnull pin_yin"`                     //拼音
	Ch                       string   `xorm:"default() notnull ch"`                          //字符
	Radical                  string   `xorm:"default() notnull radical"`                     //部首
	RadicalStroke            int      `xorm:"default(0) notnull radical_stroke"`             //部首笔画
	Stroke                   int      `xorm:"default() notnull stroke"`                      //总笔画数
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
	"部首:":     parseKangxiBuShou,
	"简体部首:":   parseSimple,
	"繁体部首:":   parseTraditional,
	"康熙字典笔画:": parseKangXi,
	"拼音":      parsePinYin,
}

func initStringArray() []string {
	return *new([]string)
}

func NewCharacter() *Character {
	return &Character{
		PinYin:                   initStringArray(),
		Ch:                       "",
		Radical:                  "",
		RadicalStroke:            0,
		IsKangXi:                 false,
		KangXi:                   "",
		KangXiStroke:             0,
		SimpleRadical:            "",
		SimpleRadicalStroke:      0,
		SimpleTotalStroke:        0,
		TraditionalRadical:       "",
		TraditionalRadicalStroke: 0,
		TraditionalTotalStroke:   0,
		NameScience:              false,
		WuXing:                   "",
		Lucky:                    "",
		Regular:                  false,
		TraditionalCharacter:     initStringArray(),
		VariantCharacter:         initStringArray(),
		Comment:                  initStringArray(),
	}
}

// Clone ...
func (c *Character) Clone() (char *Character) {
	char = new(Character)
	*char = *c
	copy(char.PinYin, c.PinYin)
	return char
}

// InsertIfNotExist ...
func (c *Character) InsertIfNotExist(session *xorm.Session) (i int64, e error) {
	i, e = session.Table(&Character{}).Where("ch = ?", c.Ch).Count()
	if e != nil {
		return
	}
	if i == 0 {
		i, e = session.InsertOne(c)
		return
	}
	i, e = session.Where("ch = ?", c.Ch).Update(c)
	return
}

func parseDummy(c *Character, index int, input string) {
	log.With("character", c, "index", index, "input", input).Warn("dummy")
}
func parseKangXi(c *Character, index int, input string) {
	if debug {
		log.With("index", index, "input", input).Info("kangxi")
	}
	switch index {
	case 2:
		input = strings.ReplaceAll(input, "(", " ")
		input = strings.ReplaceAll(input, ")", " ")
		s := strings.Split(strings.TrimSpace(input), ";")
		if len(s) > 0 {
			vv := strings.Split(strings.TrimSpace(s[0]), ":")
			if len(vv) == 2 && (vv[0] != c.Ch) {
				c.KangXi = vv[0]
				i, e := strconv.Atoi(strings.TrimSpace(vv[1]))
				if e != nil {
					log.Error(e)
					return
				}
				c.KangXiStroke = i
			} else {
				//c.KangXi = c.Ch
				//c.KangXiStroke = c.Stroke
			}
			if len(s) > 1 {
				log.Warn(s)
			}
		}

	default:
	}
}
func parseKangxiBuShou(c *Character, index int, input string) {
	if debug {
		log.With("index", index, "input", input).Info("bushou")
	}
	switch index {
	case 1:
		c.KangXi = input
		c.Radical = input
	case 3:
		parseNumber(&c.RadicalStroke, input)
	case 5:
		parseNumber(&c.KangXiStroke, input)
		parseNumber(&c.Stroke, input)
	default:
		//log.Error("bushou")
	}
}

func parseSimple(c *Character, index int, input string) {
	log.With("index", index, "input", input).Info("simple")
	switch index {
	case 1:
		parseBuShouBracket(input, &c.SimpleRadical, &c.SimpleRadicalStroke)
	case 3:
		parseNumber(&c.SimpleTotalStroke, input)
	default:
		//log.Error("simple")
	}
}

func parseTraditional(c *Character, index int, input string) {
	if debug {
		log.With("index", index, "input", input).Info("traditional")
	}
	switch index {
	case 1:
		parseBuShouBracket(input, &c.TraditionalRadical, &c.TraditionalRadicalStroke)
	case 3:
		parseNumber(&c.TraditionalTotalStroke, input)
	default:
		//log.Error("traditional")
	}
}

func parseKangXiCharacter(i int, selection *goquery.Selection, ch *Character) (e error) {
	f := parseDummy
	v := StringClearUp(selection.Find("font.colred").Contents().First().Text())
	if i == 0 {
		f = parsePinYin
	} else {
		if v, b := charList[v]; b {
			f = v
		}
	}
	if debug {
		log.With("index", i, "source", v).Info("first")
	}
	selection.Contents().Each(func(i int, selection *goquery.Selection) {
		if debug {
			log.With("text", selection.Text(), "index", selection.Index(), "num", i).Info("colred")
		}
		text := StringClearUp(selection.Text())
		f(ch, i, text)
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
	switch index {
	case 1, 5:
		if debug {
			log.With("index", index, "input", input).Info("pinyin")
		}
		input = strings.ReplaceAll(input, "[", "")
		input = strings.ReplaceAll(input, "]", "")
		parseArray(&c.PinYin, input)
	default:

	}

}
func parseBuShouBracket(input string, radical *string, stroke *int) {
	if debug {
		log.With("input", input).Info("bushou bracket")
	}
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
	if debug {
		log.With("input", input).Info("var char")
	}
	if input != "" {
		c.VariantCharacter = append(c.VariantCharacter, input)
	}
}
func parseTraditionalCharacter(c *Character, index int, input string) {
	if debug {
		log.With("input", input).Info("trad char")
	}
	if input != "" {
		c.TraditionalCharacter = append(c.TraditionalCharacter, input)
	}
}
func parseRegular(c *Character, index int, input string) {
	if debug {
		log.With("input", input).Info("regular")
	}
	if input == "是" {
		c.Regular = true
	}
}
func parseNameScience(c *Character, index int, input string) {
	if debug {
		log.With("input", input).Info("name science")
	}
	if input == "是" {
		c.NameScience = true
	}
}
func parseLucy(c *Character, index int, input string) {
	if debug {
		log.With("input", input).Info("lucky")
	}
	c.Lucky = input
}
func parseWuXing(c *Character, index int, input string) {
	if debug {
		log.With("input", input).Info("wuxing")
	}
	c.WuXing = input
}
func parseDictInformation(selection *goquery.Selection, ch *Character) (e error) {
	fn := parseDummy

	selection.Find("li").Contents().Each(func(i int, selection *goquery.Selection) {
		if debug {
			log.With("text", selection.Text(), "index", selection.Index(), "num", i).Info("li")
		}
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

	selection.Find("li").Each(func(i int, selection *goquery.Selection) {
		if debug {
			log.With("text", selection.Text(), "index", selection.Index(), "num", i).Info("li2")
		}
		tx := selection.Find("span").Text()
		if v, b := infoList2[tx]; b {
			selection.Find("a").Each(func(i int, selection *goquery.Selection) {
				v(ch, selection.Index(), selection.Text())
			})
		}
	})
	return
}
func parseComment(c *Character, index int, input string) {
	if debug {
		log.With("input", input).Info("comment")
	}
	if input == "" {
		return
	}
	c.Comment = append(c.Comment, StringClearUp(input))
}

func parseDictComment(selection *goquery.Selection, character *Character) (e error) {
	tx := StringClearUp(selection.Find("li > a").Text())
	if tx == "康熙字典解释" {
		selection.Find("li > div").Contents().Each(func(i int, selection *goquery.Selection) {
			if debug {
				log.With("text", selection.Text(), "index", selection.Index(), "num", i).Info("li3")
			}
			parseComment(character, i, selection.Text())
		})
	}

	return
}
