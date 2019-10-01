package excavator

import (
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"github.com/godcong/excavator/net"
)

const HanChengMainPage = "http://hy.httpcn.com/bushou/zi/"

func grabRadicalList(s SearchType, url string) {
	document, e := net.CacheQuery(url)
	if e != nil {
		panic(e)
	}
	//radical := make(chan *RadicalCharacter)
	rc := analyzePinyinRadical(document)
	bytes, e := json.Marshal(rc)
	if e != nil {
		return
	}
	log.With("size", len(rc)).Info(string(bytes))

	for idx := range rc {
		e := fillRadicalPinyinDetail(rc[idx])
		if e != nil {
			log.With("bushou", rc[idx].BuShou, "pinyin", rc[idx].PinYin).Error(e)
			continue
		}
	}
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

func fillRadicalPinyinDetail(character *RadicalCharacter) (err error) {
	q := NewQuery(RequestTypeOption(RequestTypeKangXiPinyin))

	closer, e := q.Grab(character.PinYin)
	if e != nil {
		return e
	}

	radical, e := RadicalReader(closer)
	if e != nil {
		return e
	}
	log.Infof("%+v", radical)
	for _, tmp := range *(*[]RadicalUnion)(radical) {
		for i := range tmp.RadicalCharacterArray {
			rc := tmp.RadicalCharacterArray[i]
			rc.BuShou = character.BuShou
			rc.Alphabet = character.Alphabet
			one, e := insertRadicalCharacter(db, &rc)
			if e != nil {
				return e
			}
			log.With("num", one).Info(rc)
		}
		log.With("value", radical).Info("radical")
	}

	return nil
}
