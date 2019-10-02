package excavator

import (
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"github.com/go-xorm/xorm"
	"github.com/godcong/excavator/net"
	"io"
	"io/ioutil"
)

type SearchType int

const (
	SearchTypePinyin SearchType = iota
	SearchTypeBushou
	SearchTypeBihua
)

// RadicalCharacter ...
type RadicalCharacter struct {
	SType    SearchType `json:"stype" json:"stype"`
	Hash     string     `json:"hash" xorm:"pk hash"`
	Zi       string     `json:"zi" xorm:"zi"`
	Alphabet string     `json:"alphabet" xorm:"alphabet"`
	PinYin   string     `json:"pinyin" xorm:"pinyin"`
	BiHua    string     `json:"bihua" xorm:"bihua"`
	BuShou   string     `json:"bushou" xorm:"bushou"`
	Num      string     `json:"num" xorm:"num"`
	URL      string     `json:"url" xorm:"url"`
}

func (r *RadicalCharacter) BeforeInsert() {
	r.Hash = net.Hash(r.URL)
}

func RadicalReader(reader io.ReadCloser) (*Radical, error) {
	bytes, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	log.Info(string(bytes))
	return UnmarshalRadical(bytes)
}

// Radical ...
type Radical []RadicalUnion

// UnmarshalRadical ...
func UnmarshalRadical(data []byte) (*Radical, error) {
	var r Radical
	err := json.Unmarshal(data, &r)
	return &r, err
}

// Marshal ...
func (r *Radical) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// RadicalUnion ...
type RadicalUnion struct {
	String                *string
	RadicalCharacterArray []RadicalCharacter
}

// UnmarshalJSON ...
func (x *RadicalUnion) UnmarshalJSON(data []byte) error {
	x.RadicalCharacterArray = nil
	_, err := unmarshalUnion(data, nil, nil, nil, &x.String, true, &x.RadicalCharacterArray, false, nil, false, nil, false, nil, false)
	if err != nil {
		return err
	}
	return nil
}

// MarshalJSON ...
func (x *RadicalUnion) MarshalJSON() ([]byte, error) {
	return marshalUnion(nil, nil, nil, x.String, x.RadicalCharacterArray != nil, x.RadicalCharacterArray, false, nil, false, nil, false, nil, false)
}

func insertOrUpdateRadicalCharacter(engine *xorm.Engine, character *RadicalCharacter) (i int64, e error) {
	tmp := &RadicalCharacter{}
	b, e := engine.Where("hash = ?", net.Hash(character.URL)).Get(tmp)
	if e != nil {
		return 0, e
	}

	log.With("url", character.URL, "character", character.BuShou, "zi", character.Zi, "hash", net.Hash(character.URL)).Info("insert")
	if !b {
		return engine.InsertOne(character)
	}
	copyRadicalCharacter(tmp, character)
	return engine.Where("hash = ?", net.Hash(character.URL)).Update(tmp)
}

func stringCompareCopy(tg *string, src string) string {
	if src != "" {
		*tg = src
	}
	return *tg
}

func copyRadicalCharacter(tg, src *RadicalCharacter) {
	stringCompareCopy(&tg.URL, src.URL)
	stringCompareCopy(&tg.Alphabet, src.Alphabet)
	stringCompareCopy(&tg.PinYin, src.PinYin)
	stringCompareCopy(&tg.BuShou, src.BuShou)
	stringCompareCopy(&tg.BiHua, src.BiHua)
	stringCompareCopy(&tg.Zi, src.Zi)
	stringCompareCopy(&tg.Num, src.Num)
}

func analyzePinyinRadical(document *goquery.Document) (rc []*RadicalCharacter) {
	document.Find("#segmentedControls > ul > li.mui-table-view-cell.mui-collapse").Each(func(i int, selection *goquery.Selection) {
		alphabet := selection.Find("a.mui-navigate-right").Text()
		selection.Find("div > a[data-action]").Each(func(i int, selection *goquery.Selection) {
			log.With("index", i, "text", selection.Text()).Info("pinyin")
			radChar := new(RadicalCharacter)
			radChar.Alphabet = alphabet
			radChar.SType = SearchTypePinyin
			radChar.PinYin, _ = selection.Attr("data-action")
			log.With("pinyin", radChar.PinYin).Info("pinyin")
			if radChar.PinYin != "" {
				rc = append(rc, radChar)
			}
		})
	})
	log.Infof("radical[%+v]", rc)
	return
}
func analyzeBushouRadical(document *goquery.Document) (rc []*RadicalCharacter) {
	document.Find("#segmentedControls > ul > li.mui-table-view-cell.mui-collapse").Each(func(i int, selection *goquery.Selection) {
		bihua := selection.Find("a.mui-navigate-right").Text()
		selection.Find("div > a[data-action]").Each(func(i int, selection *goquery.Selection) {
			log.With("index", i, "text", selection.Text()).Info("bushou")
			radChar := new(RadicalCharacter)
			radChar.BiHua = bihua
			radChar.SType = SearchTypeBushou
			radChar.BuShou, _ = selection.Attr("data-action")
			log.With("bushou", radChar.BuShou).Info("bushou")
			if radChar.BuShou != "" {
				rc = append(rc, radChar)
			}
		})
		log.Infof("radical[%+v]", rc)
	})
	return
}
