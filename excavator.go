package excavator

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/godcong/excavator/net"
	"github.com/godcong/go-trait"
	"os"
	"path/filepath"
	"strings"

	"github.com/xormsharp/xorm"
)

var log = trait.NewZapFileSugar("excavator.log")
var db *xorm.Engine
var debug = false

const tmpFile = "tmp"

// Excavator ...
type Excavator struct {
	Workspace string `json:"workspace"`
	db        *xorm.Engine
	soList    []string
	url       string
	action    []RadicalType
	//radicalType RadicalType
}

func (exc *Excavator) SoList() []string {
	return exc.soList
}

func (exc *Excavator) SetSoList(soList []string) {
	exc.soList = soList
}

// DB ...
func (exc *Excavator) DB() *xorm.Engine {
	return exc.db
}

// SetDB ...
func (exc *Excavator) SetDB(db *xorm.Engine) {
	exc.db = db
}

type ExArgs func(exc *Excavator)

func URLArgs(url string) ExArgs {
	return func(exc *Excavator) {
		exc.url = url
	}
}

func DBArgs(engine *xorm.Engine) ExArgs {
	return func(exc *Excavator) {
		exc.db = engine
	}
}

func ActionArgs(act ...RadicalType) ExArgs {
	return func(exc *Excavator) {
		exc.action = act
	}
}

// New ...
func New(args ...ExArgs) *Excavator {
	exc := &Excavator{
		Workspace: getDefaultPath(),
		url:       DefaultMainPage,
	}
	for _, arg := range args {
		arg(exc)
	}
	return exc
}

// init ...
func (exc *Excavator) init() {
	if exc.db == nil {
		exc.db = InitMysql("localhost:3306", "root", "111111")
	}
	e := exc.db.Sync2(&RadicalCharacter{}, &Character{})
	if e != nil {
		panic(e)
	}
}

// Run ...
func (exc Excavator) Run() (e error) {
	log.Info("excavator run")
	exc.init()

	for _, act := range exc.action {
		e = grabRadicalList(&exc, act)
		if e != nil {
			log.Error(e)
			panic(e)
		}

		e = parseCharacter(&exc, act)
		if e != nil {
			return e
		}

	}

	return nil
}
func fillRadicalDetail(exc *Excavator, radical *Radical, character *RadicalCharacter) (err error) {
	if debug {
		log.Infof("%+v", radical)
	}
	for _, tmp := range *(*[]RadicalUnion)(radical) {
		for i := range tmp.RadicalCharacterArray {
			rc := tmp.RadicalCharacterArray[i]
			rc.Alphabet = character.Alphabet
			rc.BiHua = character.BiHua
			rc.QiBi = character.QiBi
			rc.QBNum = character.QBNum
			rc.BHNum = character.BHNum
			rc.TotalBiHua = character.TotalBiHua
			rc.CharType = character.CharType
			one, e := insertOrUpdateRadicalCharacter(exc.db, &rc)
			if e != nil {
				return e
			}
			if debug {
				log.With("num", one).Info(rc)
			}
		}
	}
	return nil
}

func findRadical(exc *Excavator, characters chan<- *RadicalCharacter) {
	defer func() {
		characters <- nil
	}()
	i, e := exc.db.Where("char_type = ?", radicalCharType(exc.radicalType)).Count(RadicalCharacter{})
	if e != nil || i == 0 {
		log.Error(e)
		return
	}
	if debug {
		log.With("total", i).Info("total char")
	}
	for x := int64(0); x < i; x += 500 {
		rc := new([]RadicalCharacter)
		e := exc.db.Where("char_type = ?", radicalCharType(exc.radicalType)).Limit(500, int(x)).Find(rc)
		if e != nil {
			log.Error(e)
			continue
		}
		for i := range *rc {
			characters <- &(*rc)[i]
		}
	}
}

// IsExist ...
func (exc *Excavator) IsExist(name string) bool {
	_, e := os.Open(name)
	return e == nil || os.IsExist(e)
}

// GetPath ...
func getDefaultPath() string {
	wd, e := os.Getwd()
	if e != nil {
		panic(e)
	}
	return filepath.Join(wd, tmpFile)
}

/*URL 拼接地址 */
func URL(prefix string, uris ...string) string {
	end := len(prefix)
	if end > 1 && prefix[end-1] == '/' {
		prefix = prefix[:end-1]
	}

	var url = []string{prefix}
	for _, v := range uris {
		url = append(url, TrimSlash(v))
	}
	return strings.Join(url, "/")
}

// TrimSlash ...
func TrimSlash(s string) string {
	if size := len(s); size > 1 {
		if s[size-1] == '/' {
			s = s[:size-1]
		}
		if s[0] == '/' {
			s = s[1:]
		}
	}
	return s
}

func getCharacter(document *goquery.Document, c *RadicalCharacter, kangxi bool) *Character {
	ch := NewCharacter()
	ch.IsKangXi = kangxi
	ch.Ch = c.Zi
	document.Find("div.hanyu-tujie.mui-clearfix > div.info > p.mui-ellipsis").Each(func(i int, selection *goquery.Selection) {
		if ch.IsKangXi {
			ch.KangXi = c.Zi
			e := parseKangXiCharacter(i, selection, ch)
			if e != nil {
				log.Error(e)
			}
		} else {
			e := parseZiCharacter(i, selection, ch)
			if e != nil {
				log.Error(e)
			}
		}
	})

	document.Find("div > ul.hanyu-cha-info.mui-clearfix").Each(func(i int, selection *goquery.Selection) {
		e := parseDictInformation(selection, ch)
		if debug {
			log.Infof("%+v", ch)
		}
		if e != nil {
			log.Error(e)
		}
	})
	document.Find("div > ul.hanyu-cha-ul").Each(func(i int, selection *goquery.Selection) {
		e := parseDictComment(selection, ch)
		if e != nil {
			log.Error(e)
		}
	})
	return ch
}

func parseCharacter(exc *Excavator, radicalType RadicalType) (e error) {
	ch := make(chan *RadicalCharacter)
	go findRadical(exc, ch)
ParseEnd:
	for {
		select {
		case c := <-ch:
			if c == nil {
				break ParseEnd
			}
			log.With("url", c.URL).Info("character")
			document, e := net.CacheQuery(characterURL(exc.url, radicalType, c.URL))
			if e != nil {
				log.Error(e)
				continue
			}
			character := getCharacter(document, c, isKangxi(radicalType))
			_, e = character.InsertOrUpdate(exc.db.Where(""))
			if e != nil {
				return e
			}
		}
	}
	return nil
}

func characterURL(m string, rt RadicalType, url string) string {
	switch rt {
	case RadicalTypeKangXiPinyin, RadicalTypeKangXiBushou, RadicalTypeKangXiBihua:
		return URL(m, "html/kangxi", url)
	}
	return URL(m, "html/zi", url)
}
func radicalCharType(radicalType RadicalType) string {
	switch radicalType {
	case RadicalTypeHanChengPinyin, RadicalTypeHanChengBushou, RadicalTypeHanChengBihua:
		return "hancheng"
	default:
		//RadicalTypeKangXiBihua, RadicalTypeKangXiPinyin, RadicalTypeKangXiBushou:
		return "kangxi"
	}
	//return ""
}
