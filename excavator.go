package excavator

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/godcong/excavator/net"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-xorm/xorm"
	"github.com/godcong/go-trait"
)

var log = trait.NewZapFileSugar("excavator.log")
var db *xorm.Engine
var debug = false

const tmpFile = "tmp"

// Excavator ...
type Excavator struct {
	Workspace   string `json:"workspace"`
	db          *xorm.Engine
	url         string
	action      []RadicalType
	radicalType RadicalType
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

func ActionArgs(act ...RadicalType) ExArgs {
	return func(exc *Excavator) {
		exc.action = act
	}
}

// New ...
func New(radicalType RadicalType, args ...ExArgs) *Excavator {
	exc := &Excavator{
		radicalType: radicalType,
		Workspace:   getDefaultPath(),
		url:         DefaultMainPage,
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
func (exc *Excavator) Run() error {
	log.Info("excavator run")
	exc.init()

	for _, act := range exc.action {
		excClone := *exc
		excClone.radicalType = act
		e := grabRadicalList(&excClone)
		if e != nil {
			log.Error(e)
			panic(e)
		}
	}

	e := parseCharacter(exc)
	if e != nil {
		return e
	}
	return nil
}
func fillRadicalDetail(exc *Excavator, radical *Radical, character *RadicalCharacter) (err error) {
	log.Infof("%+v", radical)
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
			log.With("num", one).Info(rc)
		}
		log.With("value", radical).Info("radical")
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
	log.With("total", i).Info("total char")
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

func getCharacter(document *goquery.Document) *Character {
	ch := NewCharacter()
	document.Find("div.hanyu-tujie.mui-clearfix > div.info > p.mui-ellipsis").Each(func(i int, selection *goquery.Selection) {
		e := parseKangXiCharacter(i, selection, ch)
		if e != nil {
			log.Error(e)
		}
	})

	document.Find("div > ul.hanyu-cha-info.mui-clearfix").Each(func(i int, selection *goquery.Selection) {
		e := parseDictInformation(selection, ch)
		log.Infof("%+v", ch)
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

func parseCharacter(exc *Excavator) (e error) {
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
			document, e := net.CacheQuery(characterURL(exc, c.URL))
			if e != nil {
				log.Error(e)
				continue
			}
			character := getCharacter(document)
			character.Ch = c.Zi
			_, e = character.InsertIfNotExist(exc.db.Where(""))
			if e != nil {
				return e
			}
		}
	}
	return nil
}

func characterURL(excavator *Excavator, url string) string {
	switch excavator.radicalType {
	case RadicalTypeKangXiPinyin, RadicalTypeKangXiBushou, RadicalTypeKangXiBihua:
		return URL(excavator.url, "html/kangxi", url)
	default:
		return URL(excavator.url, "html/zi", url)
	}
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
