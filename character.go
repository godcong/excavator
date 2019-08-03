package excavator

import (
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

// CharacterFunc ...
type CharacterFunc func(character *Character) error

// ClassRegular ...
const ClassRegular int = iota

// RootRadicalCharacter ...
type RootRadicalCharacter struct {
	Class     int      `json:"class"`     //分类
	Character string   `json:"character"` //字符
	Link      string   `json:"link"`      //链接
	Pinyin    []string `json:"pinyin"`    //拼音
}

//BasicExplanation 基本解释
type BasicExplanation struct {
	Pinyin       string   `json:"pinyin"`   //拼音
	Phonetic     string   `json:"phonetic"` //注音
	BasicMeaning []string `xorm:"json basic_meaning" json:"basic_meaning"`
	OtherMeaning []string
}

//DetailedExplanation  详细解释
type DetailedExplanation struct {
	Character       string   `json:"character"`
	Pinyin          string   `json:"pinyin"`
	DetailedMeaning []string `json:"detailed_meaning"`
}

// MandarinDictionary 国语辞典
type MandarinDictionary struct {
	PartOfSpeech string   //词性
	Pinyin       string   //拼音
	Phonetic     string   //注音
	Explanation  []string //解释
}

//KangxiDictionary 康熙字典
type KangxiDictionary struct {
}

//StandardCharacter 标准字符
type StandardCharacter struct {
	Radical             string
	BasicExplanation    BasicExplanation
	DetailedExplanation []DetailedExplanation
	MandarinDictionary  []MandarinDictionary
	KangxiDictionary    KangxiDictionary
	CharacterDetail     map[string]string
}

// BaseCharacter ...
type BaseCharacter struct {
	NeedFix   bool
	Character string
	Data      map[string]string
}

//Character 字符
type Character struct {
	PinYin                   []string `xorm:"pin_yin"`                    //拼音
	Character                string   `xorm:"character"`                  //字符
	Radical                  string   `xorm:"radical"`                    //部首
	RadicalStroke            int      `xorm:"radical_stroke"`             //部首笔画
	KangXi                   string   `json:"traditional_radical"`        //康熙
	KangXiStroke             int      `json:"traditional_radical_stroke"` //康熙笔画
	SimpleRadical            string   `json:"traditional_radical"`        //简体部首
	SimpleRadicalStroke      int      `json:"traditional_radical_stroke"` //简体部首笔画
	SimpleTotalStroke        int      `json:"traditional_radical_stroke"` //简体笔画
	TraditionalRadical       string   `json:"traditional_radical"`        //繁体部首
	TraditionalRadicalStroke int      `json:"traditional_radical_stroke"` //繁体部首笔画
	TraditionalTotalStroke   int      `json:"traditional_radical_stroke"` //简体部首笔画
	NameScience              bool     `json:"name_science"`               //姓名学
	WuXing                   string   `json:"wu_xing"`                    //五行
	Lucky                    string   `json:"lucky"`                      //吉凶寓意
	Regular                  bool     `json:"regular"`                    //常用
	TraditionalCharacter     []string `json:"traditional_character"`      //繁体字
	VariantCharacter         []string `json:"variant_character"`          //异体字
}

//Folk 民俗参考
type Folk struct {
	CommonlyCharacters   string //是否为常用字
	NameScience          string //姓名学
	FiveElementCharacter string //汉字五行
	GodBadMoral          string //吉凶寓意
}

//Structure 字形结构
type Structure struct {
	DecompositionSearch string //首尾分解查字
	StrokeNumber        string //笔顺编号
	StrokeReadWrite     string //笔顺读写
}

//Explain 说文解字
type Explain struct {
	Intro  string //简介
	Detail string //详情
}

//PhoneticDialect 音韵方言
type PhoneticDialect struct {
	GuangYun  string //广　韵
	Mandarin  string //国　语
	Cantonese string //粤　语
}

//Index 索引参考
type Index struct {
	AncientWrite      string //古文字诂林
	HometownTrain     string //故训彙纂
	Explain           string //说文解字
	KangxiDictionary  string //康熙字典
	ChineseDictionary string //汉语字典
	Cihai             string //辞　海
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
	"繁体字集：":   parseTraditionalCharacter,
	"异体字集：":   parseVariantCharacter,
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

	n.Find("li").Each(func(i int, selection *goquery.Selection) {
		log.With("text", selection.Text(), "index", selection.Index(), "num", i).Info("li")
		tx := selection.Text()
		if i == 0 {
			fn = parseDummy
			if v, b := infoList[tx]; b {
				fn = v
			}
		}

		if goquery.NodeName(selection) == "#text" {
			fn(ch, i, StringClearUp(tx))
		}
	})

	return
}
