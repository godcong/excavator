package excavator

import (
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"github.com/go-xorm/xorm"
	"github.com/godcong/excavator/net"
	"io/ioutil"
)

type RadicalType int

const DefaultMainPage = "http://hy.httpcn.com"

const (
	RadicalTypeHanChengPinyin RadicalType = iota
	RadicalTypeHanChengBushou
	RadicalTypeHanChengBihua
	RadicalTypeKangXiPinyin
	RadicalTypeKangXiBushou
	RadicalTypeKangXiBihua
)

// RadicalCharacter ...
type RadicalCharacter struct {
	Hash       string `json:"hash" xorm:"pk hash"`
	CharType   string `json:"char_type" json:"char_type"`
	Zi         string `json:"zi" xorm:"zi"`
	Alphabet   string `json:"alphabet" xorm:"alphabet"`
	PinYin     string `json:"pinyin" xorm:"pinyin"`
	BiHua      string `json:"bihua" xorm:"bihua"`
	BuShou     string `json:"bushou" xorm:"bushou"`
	TotalBiHua string `json:"total_bihua" xorm:"total_bihua"`
	QiBi       string `json:"qibi" xorm:"qibi"`
	BHNum      string `json:"bh_num" xorm:"bh_num"`
	QBNum      string `json:"qb_num" xorm:"qb_num"`
	Num        string `json:"num" xorm:"num"`
	URL        string `json:"url" xorm:"url"`
}

func (r *RadicalCharacter) BeforeInsert() {
	r.Hash = r.GenHash()
}

func (r *RadicalCharacter) GenHash() string {
	return net.Hash(r.CharType + "_" + r.Zi)
}

func RadicalReader(radicalType RadicalType, wd string, qb string) (*Radical, error) {
	reader, e := NewQuery().Grab(radicalType)(wd, qb)
	if e != nil {
		return nil, e
	}

	bytes, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	if debug {
		log.With("info", string(bytes)).Info("radical reader")
	}
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
	b, e := engine.Where("hash = ?", character.GenHash()).Get(tmp)
	if e != nil {
		return 0, e
	}
	if debug {
		log.With("url", character.URL, "character", character.BuShou, "zi", character.Zi, "hash", net.Hash(character.Zi)).Info("insert")
	}
	if !b {
		return engine.InsertOne(character)
	}
	copyRadicalCharacter(tmp, character)
	return engine.Where("hash = ?", character.GenHash()).Update(tmp)
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
	stringCompareCopy(&tg.TotalBiHua, src.TotalBiHua)
	stringCompareCopy(&tg.QBNum, src.QBNum)
	stringCompareCopy(&tg.BHNum, src.BHNum)
	stringCompareCopy(&tg.QiBi, src.QiBi)
	stringCompareCopy(&tg.Zi, src.Zi)
	stringCompareCopy(&tg.Num, src.Num)
	stringCompareCopy(&tg.CharType, src.CharType)
}

func analyzePinyinRadical(document *goquery.Document) (rc []*RadicalCharacter) {
	document.Find("#segmentedControls > ul > li.mui-table-view-cell.mui-collapse").Each(func(i int, selection *goquery.Selection) {
		alphabet := selection.Find("a.mui-navigate-right").Text()
		log.Info(selection.Html())
		selection.Find("div > a[data-action]").Each(func(i int, selection *goquery.Selection) {
			log.With("index", i, "text", selection.Text()).Info("pinyin")
			radChar := new(RadicalCharacter)
			radChar.Alphabet = alphabet
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
		log.Info(selection.Html())
		selection.Find("div > a[data-action]").Each(func(i int, selection *goquery.Selection) {
			log.With("index", i, "text", selection.Text()).Info("bushou")
			radChar := new(RadicalCharacter)
			radChar.BiHua = bihua
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

func analyzeBihuaRadical(document *goquery.Document) (rc []*RadicalCharacter) {
	document.Find("#segmentedControls > ul > li.mui-table-view-cell.mui-collapse").Each(func(i int, selection *goquery.Selection) {
		tbihua := selection.Find("a.mui-navigate-right").Text()
		log.Info(selection.Html())
		selection.Find("div > a[data-bh]").Each(func(i int, selection *goquery.Selection) {
			log.With("index", i, "text", selection.Text()).Info("bihua")
			radChar := new(RadicalCharacter)
			radChar.TotalBiHua = tbihua
			radChar.QBNum, _ = selection.Attr("data-val")
			radChar.BHNum, _ = selection.Attr("data-bh")
			radChar.QiBi, _ = selection.Attr("data-qb")
			log.With("qibi", radChar.QiBi, "bh", radChar.BHNum, "qb", radChar.QBNum).Info("bihua")
			if radChar.QiBi != "" {
				rc = append(rc, radChar)
			}
		})
		log.Infof("radical[%+v]", rc)
	})
	return
}

func getMainURL(radicalType RadicalType, url string) string {
	switch radicalType {
	case RadicalTypeHanChengPinyin:
		return url + HanChengPinyin
	case RadicalTypeHanChengBushou:
		return url + HanChengBushou
	case RadicalTypeHanChengBihua:
		return url + HanChengBihua
	case RadicalTypeKangXiBihua:
		return url + KangXiBihua
	case RadicalTypeKangXiPinyin:
		return url + KangXiPinyin
	case RadicalTypeKangXiBushou:
		return url + KangXiBushou
	}
	return ""
}

func grabRadicalList(exc *Excavator) (e error) {
	document, e := net.CacheQuery(getMainURL(exc.radicalType, exc.url))
	if e != nil {
		return e
	}
	var rc []*RadicalCharacter
	switch exc.radicalType {
	case RadicalTypeHanChengPinyin:
		rc = analyzePinyinRadical(document)
		bytes, e := json.Marshal(rc)
		if e != nil {
			return e
		}
		log.With("size", len(rc)).Info(string(bytes))
		for idx := range rc {
			radical, e := RadicalReader(exc.radicalType, rc[idx].PinYin, "")
			if e != nil {
				return e
			}
			char := rc[idx]
			char.CharType = "hancheng"
			e = fillRadicalDetail(exc, radical, char)
			if e != nil {
				log.With("bushou", rc[idx].BuShou, "pinyin", rc[idx].PinYin).Error(e)
				continue
			}
		}
	case RadicalTypeHanChengBushou:
		rc = analyzeBushouRadical(document)
		bytes, e := json.Marshal(rc)
		if e != nil {
			return e
		}
		log.With("size", len(rc)).Info(string(bytes))
		for idx := range rc {
			radical, e := RadicalReader(exc.radicalType, rc[idx].BuShou, "")
			if e != nil {
				return e
			}
			char := rc[idx]
			char.CharType = "hancheng"
			e = fillRadicalDetail(exc, radical, char)
			if e != nil {
				log.With("bushou", rc[idx].BuShou, "pinyin", rc[idx].PinYin).Error(e)
				continue
			}
		}
	case RadicalTypeHanChengBihua:
		rc = analyzeBihuaRadical(document)
		bytes, e := json.Marshal(rc)
		if e != nil {
			return e
		}
		log.With("size", len(rc)).Info(string(bytes))
		for idx := range rc {
			radical, e := RadicalReader(exc.radicalType, rc[idx].BHNum, rc[idx].QBNum)
			if e != nil {
				return e
			}
			char := rc[idx]
			char.CharType = "hancheng"
			e = fillRadicalDetail(exc, radical, char)
			if e != nil {
				log.With("bushou", rc[idx].BuShou, "pinyin", rc[idx].PinYin, "bihua", rc[idx].BiHua).Error(e)
				continue
			}
		}
	case RadicalTypeKangXiBushou:
		rc = analyzeBushouRadical(document)
		bytes, e := json.Marshal(rc)
		if e != nil {
			return e
		}
		log.With("size", len(rc)).Info(string(bytes))
		for idx := range rc {
			radical, e := RadicalReader(exc.radicalType, rc[idx].BuShou, "")
			if e != nil {
				return e
			}
			char := rc[idx]
			char.CharType = "kangxi"
			e = fillRadicalDetail(exc, radical, char)
			if e != nil {
				log.With("bushou", rc[idx].BuShou, "pinyin", rc[idx].PinYin, "bihua", rc[idx].BiHua).Error(e)
				continue
			}
		}
	case RadicalTypeKangXiPinyin:
		rc = analyzePinyinRadical(document)
		bytes, e := json.Marshal(rc)
		if e != nil {
			return e
		}
		log.With("size", len(rc)).Info(string(bytes))
		for idx := range rc {
			radical, e := RadicalReader(exc.radicalType, rc[idx].PinYin, "")
			if e != nil {
				return e
			}
			char := rc[idx]
			char.CharType = "kangxi"
			e = fillRadicalDetail(exc, radical, char)
			if e != nil {
				log.With("bushou", rc[idx].BuShou, "pinyin", rc[idx].PinYin, "bihua", rc[idx].BiHua).Error(e)
				continue
			}
		}
	case RadicalTypeKangXiBihua:
		rc = analyzeBihuaRadical(document)
		bytes, e := json.Marshal(rc)
		if e != nil {
			return e
		}
		log.With("size", len(rc)).Info(string(bytes))
		for idx := range rc {
			radical, e := RadicalReader(exc.radicalType, rc[idx].BHNum, rc[idx].QBNum)
			if e != nil {
				return e
			}
			char := rc[idx]
			char.CharType = "kangxi"
			e = fillRadicalDetail(exc, radical, char)
			if e != nil {
				log.With("bushou", rc[idx].BuShou, "pinyin", rc[idx].PinYin, "bihua", rc[idx].BiHua).Error(e)
				continue
			}
		}
	}
	return nil
}
