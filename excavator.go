package excavator

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/antchfx/htmlquery"
	"github.com/free-utils-go/cachenet"
	"github.com/godcong/excavator/config"
	"github.com/godcong/excavator/models"
	"github.com/godcong/go-trait"
	"golang.org/x/net/html"
	"xorm.io/xorm"
)

var Log = trait.NewZapFileSugar("excavator.log")

var debug = false

// Excavator ...
type Excavator struct {
	Db           *xorm.Engine
	DbFate       *xorm.Engine
	base_url     string
	unicode_file string
	cache        *cachenet.Cache
	action       string
}

// New ...
func New(args ...func(exc *Excavator)) *Excavator {
	cfg := config.LoadConfig()

	exc := &Excavator{
		Db:           InitXorm(&cfg.DatabaseExc),
		DbFate:       InitXorm(&cfg.DatabaseFate),
		base_url:     cfg.BaseUrl,
		unicode_file: cfg.UnicodeFile,
		cache:        cachenet.NewCache(cfg.TmpDir),
	}

	for _, arg := range args {
		arg(exc)
	}
	return exc
}

// Run ...
func (exc *Excavator) Run() (e error) {
	if exc.Db == nil {
		panic("没有初始化")
	}

	Log.Info("excavator run")

	switch exc.action {
	case config.ActionGrab:
		e = grabHanChengList(exc)
		if e != nil {
			Log.Error(e)
			panic(e)
		}

		ResetExc(exc.Db)
		ResetFate(exc.DbFate)

		e = parseCharacter(exc)
		if e != nil {
			return e
		}

		e = variantChars(exc)
		if e != nil {
			return e
		}

		e = simplifyChars(exc)
		if e != nil {
			return e
		}
	case config.ActionParse:
		ResetExc(exc.Db)
		ResetFate(exc.DbFate)

		e = parseCharacter(exc)
		if e != nil {
			return e
		}

		e = variantChars(exc)
		if e != nil {
			return e
		}

		e = simplifyChars(exc)
		if e != nil {
			return e
		}
	case config.ActionVariant:
		ResetFate(exc.DbFate)

		e = variantChars(exc)
		if e != nil {
			return e
		}

		e = simplifyChars(exc)
		if e != nil {
			return e
		}
	case config.ActionSimplify:
		ResetFate(exc.DbFate)

		e = simplifyChars(exc)
		if e != nil {
			return e
		}
	default:
		fmt.Printf("设置config.json中的action。有五种选择“%s”，“%s”，“%s”，“%s”，“”\n", config.ActionSimplify,
			config.ActionVariant, config.ActionParse, config.ActionGrab)
	}

	return nil
}

//从数据库取出汉程字符链接
func getHanChengFromDB(exc *Excavator, hanCheng chan<- *models.HanChengChar) {
	defer func() {
		hanCheng <- nil
	}()
	i, e := exc.Db.Count(models.HanChengChar{})
	if e != nil || i == 0 {
		Log.Error(e)
		return
	}
	if debug {
		Log.With("total", i).Info("total char")
	}
	for x := 0; x < int(i); x += 500 {
		rc := new([]models.HanChengChar)
		e := exc.Db.Limit(500, x).Find(rc)
		if e != nil {
			Log.Error(e)
			continue
		}
		for i := range *rc {
			hanCheng <- &(*rc)[i]
		}
	}
}

//解析汉程桌面版的文字信息
func getCharacter(exc *Excavator, unid rune, html_node *html.Node) (err error) {
	he_xin_block := htmlquery.FindOne(html_node, "//p[contains(@class, 'text15')]")

	if he_xin_block == nil {
		return errors.New("核心区找不到")
	}

	//基本解释
	ji_ben_block := htmlquery.FindOne(html_node, "//div[contains(@class, 'content16')]/span[contains(text(), '基')]/..")

	if ji_ben_block == nil {
		panic("基本解释区找不到")
	}

	err = parseComment(exc, unid, html_node, ji_ben_block)

	if err == nil {
		parseJiBenVariant(exc, unid, html_node, ji_ben_block)
	}

	parseVariant(exc, unid, html_node, he_xin_block, ji_ben_block)

	kang := parseKangXiStroke(unid, he_xin_block)

	glyph := models.Glyph{
		Unid: unid,
	}

	parseZiCharacter(exc, &glyph, unid, html_node, he_xin_block, ji_ben_block)

	parseKangXi(exc, unid, &glyph, html_node, kang)

	parseBianMa(exc, unid, html_node)

	parseDictInformation(exc, unid, html_node)

	parseYinYun(exc, unid, html_node, he_xin_block)

	parseSuoYin(exc, unid, html_node)

	parseXiangXi(exc, unid, html_node)

	parseHanYuDaZiDian(exc, unid, html_node)

	shuo_wen_block := htmlquery.FindOne(html_node, "//div[@id='div_a5']")

	parseShuoWenJieZi(exc, unid, html_node, shuo_wen_block)

	parseYanBian(exc, unid, html_node, shuo_wen_block)

	cis := map[string]bool{}

	parseChengYu(exc, unid, html_node, cis)

	parseShiCi(exc, unid, html_node, cis)

	parseCiYu(exc, unid, html_node, cis)

	return nil
}

//从汉程解析每一个字
func parseCharacter(exc *Excavator) (err error) {
	hcc := make(chan *models.HanChengChar, 100)
	go getHanChengFromDB(exc, hcc)

	invalids := map[rune]bool{}

	for {
		c := <-hcc
		if c == nil {
			break
		}
		Log.With("url", c.Url).Info("character")

		html_node, e := cachenet.CacheQuery(cachenet.UrlMerge(exc.base_url, c.Url))

		if e != nil {
			Log.Error(e)
			continue
		}

		e = getCharacter(exc, c.Unid, html_node)
		if e != nil {
			if e.Error() == "核心区找不到" || e.Error() == "基本解释同义字格式不对" {
				invalids[c.Unid] = true
				continue
			} else if e.Error() == "没有部首信息" {
				invalids[c.Unid] = true
				continue
			}
			return e
		}
	}

	for k := range invalids {
		fmt.Println(strconv.FormatUint(uint64(k), 16), " ", string(rune(k)))
	}

	fmt.Println(len(invalids))

	return nil
}
