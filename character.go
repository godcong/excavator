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
	Character                string   `xorm:"character"`                  //字符
	Radical                  string   `xorm:"radical"`                    //部首
	RadicalStroke            int      `xorm:"radical_stroke"`             //部首笔画
	SimpleRadical            string   `json:"traditional_radical"`        //简体部首
	SimpleRadicalStroke      int      `json:"traditional_radical_stroke"` //简体部首笔画
	TraditionalRadical       string   `json:"traditional_radical"`        //繁体部首
	TraditionalRadicalStroke int      `json:"traditional_radical_stroke"` //繁体部首笔画
	KangXi                   string   `json:"traditional_radical"`        //康熙
	KangXiStroke             int      `json:"traditional_radical_stroke"` //康熙笔画
	TotalStroke              int      `xorm:"total_stroke"`               //总笔画
	PinYin                   []string `xorm:"pin_yin"`                    //拼音
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
type ParseFunc func(*Character, string)

var charList = map[string]ParseFunc{
	"部首:":   parseBuShou,
	"部首笔画:": parseBuShouStroke,
	"总笔画:":  parseTotalStroke,
	"拼音":    parsePinYin,
}

func parseDummy(c *Character, input string) {
	log.With("character", c, "input", input).Info("dummy")
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
func parsePinYin(c *Character, input string) {
	log.With("input", input).Info("pinyin")
	parseArray(&c.PinYin, input)
}
func parseBuShou(c *Character, input string) {
	log.With("input", input).Info("bushou")
	c.Radical = input
}
func parseBuShouStroke(c *Character, input string) {
	parseNumber(&c.RadicalStroke, input)
}
func parseSimpleRadical(c *Character, input string) {
	log.With("input", input).Info("simple radical")
	c.SimpleRadical = input
}
func parseSimpleRadicalStroke(c *Character, input string) {
	parseNumber(&c.SimpleRadicalStroke, input)
}
func parseTraditionalRadical(c *Character, input string) {
	log.With("input", input).Info("simple radical")
	c.TraditionalRadical = input
}
func parseTraditionalRadicalStroke(c *Character, input string) {
	parseNumber(&c.TraditionalRadicalStroke, input)
}
func parseTotalStroke(c *Character, input string) {
	parseNumber(&c.TotalStroke, input)
}

func parseCharacter(element *colly.HTMLElement, ch *Character) (e error) {
	html, e := element.DOM.Html()
	if e != nil {
		log.Error(e)
		return e
	}
	n, e := goquery.NewDocumentFromReader(strings.NewReader(html))
	if e != nil {
		log.Error(e)
		return e
	}
	v := StringClearUp(n.ReplaceWith("font[class=colred]").Text())
	vv := strings.Split(v, " ")
	log.With("source", v).Info(len(vv), ":", vv)
	n1, e := goquery.NewDocumentFromReader(strings.NewReader(html))
	if e != nil {
		log.Error(e)
		return nil
	}
	n1.Find("font[class=colred]").Each(func(i int, selection *goquery.Selection) {
		log.With("text", selection.Text(), "index", selection.Index(), "num", i).Info("colred")
		text := StringClearUp(selection.Text())
		f := parseDummy
		if v, b := charList[text]; b {
			f = v
		}
		if len(vv) > i {
			f(ch, vv[i])
		}
	})
	return nil
}
