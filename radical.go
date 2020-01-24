package excavator

import (
	"encoding/json"
	"errors"
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
	RadicalTypeHanChengSo
	RadicalTypeKangXiPinyin
	RadicalTypeKangXiBushou
	RadicalTypeKangXiBihua
	RadicalTypeKangXiSo
)

type SoCharacter struct {
	Zi     string `json:"zi"`
	URL    string `json:"url"`
	Py     string `json:"py"`
	Bushou string `json:"bushou"`
	Num    string `json:"num"`
}

type SoCharacterElement struct {
	Integer     *int64
	SoCharacter *SoCharacter
}

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

// RadicalUnion ...
type RadicalUnion struct {
	String                *string
	RadicalCharacterArray []RadicalCharacter
}

// Radical ...
type Radical []RadicalUnion

type RadicalSo [][]RadicalSoElement

type RadicalSoClass struct {
	Zi     string `json:"zi"`
	URL    string `json:"url"`
	Py     string `json:"py"`
	Bushou string `json:"bushou"`
	Num    string `json:"num"`
}

type RadicalSoElement struct {
	Integer        *int64
	RadicalSoClass *RadicalSoClass
}

func (x *RadicalSoElement) UnmarshalJSON(data []byte) error {
	x.RadicalSoClass = nil
	var c RadicalSoClass
	object, err := unmarshalUnion(data, &x.Integer, nil, nil, nil, false, nil, true, &c, false, nil, false, nil, false)
	if err != nil {
		return err
	}
	if object {
		x.RadicalSoClass = &c
	}
	return nil
}

func (x *RadicalSoElement) MarshalJSON() ([]byte, error) {
	return marshalUnion(x.Integer, nil, nil, nil, false, nil, x.RadicalSoClass != nil, x.RadicalSoClass, false, nil, false, nil, false)
}

// UnmarshalRadical ...
func UnmarshalRadical(data []byte) (*Radical, error) {
	var r Radical
	err := json.Unmarshal(data, &r)
	return &r, err
}

func UnmarshalRadicalSo(data []byte) (*RadicalSo, error) {
	var r RadicalSo
	err := json.Unmarshal(data, &r)
	return &r, err
}

func (r *RadicalSo) Marshal() ([]byte, error) {
	return json.Marshal(r)
}
func (so *RadicalSo) Radical() *Radical {
	if so == nil || len(([][]RadicalSoElement)(*so)) < 2 || len(([][]RadicalSoElement)(*so)[1]) == 0 {
		return nil
	}

	elements := ([][]RadicalSoElement)(*so)[1]

	var rs []RadicalCharacter

	for _, s := range elements {
		rs = append(rs, RadicalCharacter{
			Hash:       "",
			CharType:   "",
			Zi:         s.RadicalSoClass.Zi,
			Alphabet:   "",
			PinYin:     s.RadicalSoClass.Py,
			BiHua:      "",
			BuShou:     s.RadicalSoClass.Bushou,
			TotalBiHua: "",
			QiBi:       "",
			BHNum:      "",
			QBNum:      "",
			Num:        s.RadicalSoClass.Num,
			URL:        s.RadicalSoClass.URL,
		})
	}

	return &Radical{
		RadicalUnion{
			String:                nil,
			RadicalCharacterArray: rs,
		},
	}

}

func (r *RadicalCharacter) BeforeInsert() {
	r.Hash = r.GenHash()
}

func (r *RadicalCharacter) GenHash() string {
	return net.Hash(r.CharType + "_" + r.Zi)
}

func RadicalReaderSo(radicalType RadicalType, wd string, qb string) (*RadicalSo, error) {
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
	return UnmarshalRadicalSo(bytes)
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

// Marshal ...
func (r *Radical) Marshal() ([]byte, error) {
	return json.Marshal(r)
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
		if debug {
			log.Info(selection.Html())
		}
		selection.Find("div > a[data-action]").Each(func(i int, selection *goquery.Selection) {
			if debug {
				log.With("index", i, "text", selection.Text()).Info("pinyin")
			}
			radChar := new(RadicalCharacter)
			radChar.Alphabet = alphabet
			radChar.PinYin, _ = selection.Attr("data-action")
			if debug {
				log.With("pinyin", radChar.PinYin).Info("pinyin")
			}
			if radChar.PinYin != "" {
				rc = append(rc, radChar)
			}
		})
	})
	if debug {
		log.Infof("radical[%+v]", rc)
	}
	return
}
func analyzeBushouRadical(document *goquery.Document) (rc []*RadicalCharacter) {
	document.Find("#segmentedControls > ul > li.mui-table-view-cell.mui-collapse").Each(func(i int, selection *goquery.Selection) {
		bihua := selection.Find("a.mui-navigate-right").Text()
		if debug {
			log.Info(selection.Html())
		}
		selection.Find("div > a[data-action]").Each(func(i int, selection *goquery.Selection) {
			if debug {
				log.With("index", i, "text", selection.Text()).Info("bushou")
			}
			radChar := new(RadicalCharacter)
			radChar.BiHua = bihua
			radChar.BuShou, _ = selection.Attr("data-action")
			if debug {
				log.With("bushou", radChar.BuShou).Info("bushou")
			}
			if radChar.BuShou != "" {
				rc = append(rc, radChar)
			}
		})
		if debug {
			log.Infof("radical[%+v]", rc)
		}
	})
	return
}

func analyzeBihuaRadical(document *goquery.Document) (rc []*RadicalCharacter) {
	document.Find("#segmentedControls > ul > li.mui-table-view-cell.mui-collapse").Each(func(i int, selection *goquery.Selection) {
		tbihua := selection.Find("a.mui-navigate-right").Text()
		if debug {
			log.Info(selection.Html())
		}
		selection.Find("div > a[data-bh]").Each(func(i int, selection *goquery.Selection) {
			if debug {
				log.With("index", i, "text", selection.Text()).Info("bihua")
			}
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
		if debug {
			log.Infof("radical[%+v]", rc)
		}
	})
	return
}

func isKangxi(radicalType RadicalType) bool {
	switch radicalType {
	case RadicalTypeKangXiBihua,
		RadicalTypeKangXiPinyin,
		RadicalTypeKangXiBushou:
		return true
	}
	return false
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
		//case RadicalTypeKangXiSo:
		//	return url + "/so/kangxi/"
		//case RadicalTypeHanChengSo:
		//	return url + "/so"
	}
	return ""
}

func grabRadicalList(exc *Excavator) (e error) {
	mainURL := getMainURL(exc.radicalType, exc.url)
	if mainURL == "" && (exc.radicalType != RadicalTypeKangXiSo && exc.radicalType != RadicalTypeHanChengSo) {
		return errors.New("wrong type")
	}
	var document *goquery.Document
	if mainURL != "" {
		document, e = net.CacheQuery(mainURL)
		if e != nil {
			return e
		}
	}

	var rc []*RadicalCharacter
	switch exc.radicalType {
	case RadicalTypeHanChengPinyin:
		rc = analyzePinyinRadical(document)
		if debug {
			bytes, e := json.Marshal(rc)
			if e == nil {
				log.With("size", len(rc)).Info(string(bytes))
			}
		}
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
		if debug {
			log.With("size", len(rc)).Info(string(bytes))
		}
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
		if debug {
			log.With("size", len(rc)).Info(string(bytes))
		}
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
		if debug {
			log.With("size", len(rc)).Info(string(bytes))
		}
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
		if debug {
			log.With("size", len(rc)).Info(string(bytes))
		}
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
		if debug {
			log.With("size", len(rc)).Info(string(bytes))
		}
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
	case RadicalTypeKangXiSo:
		for _, wd := range exc.SoList() {
			radical, e := RadicalReaderSo(exc.radicalType, wd, "")
			if e != nil {
				return e
			}
			char := &RadicalCharacter{
				Zi: wd,
			}
			char.CharType = "kangxi"
			e = fillRadicalDetail(exc, radical.Radical(), char)
			if e != nil {
				log.With("so", wd).Error(e)
				continue
			}
		}
	case RadicalTypeHanChengSo:
		for _, wd := range exc.SoList() {
			radical, e := RadicalReaderSo(exc.radicalType, wd, "")
			if e != nil {
				return e
			}
			char := &RadicalCharacter{
				Zi: wd,
			}
			char.CharType = "kangxi"
			e = fillRadicalDetail(exc, radical.Radical(), char)
			if e != nil {
				log.With("so", wd).Error(e)
				continue
			}
		}
	}
	return nil
}
