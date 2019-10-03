package excavator

import (
	"encoding/json"
	"github.com/godcong/excavator/net"
)

const HanChengBushou = "/bushou/zi/"
const HanChengPinyin = "/pinyin/zi/"
const HanChengBihua = "/bihua/zi/"

func grabRadicalList(s RadicalType, url string) {
	document, e := net.CacheQuery(url)
	if e != nil {
		panic(e)
	}
	var rc []*RadicalCharacter
	switch s {
	case RadicalTypeHanChengPinyin:
		rc = analyzePinyinRadical(document)
		bytes, e := json.Marshal(rc)
		if e != nil {
			return
		}
		log.With("size", len(rc)).Info(string(bytes))
		for idx := range rc {
			radical, e := RadicalReader(s, rc[idx].PinYin)
			if e != nil {
				return
			}
			char := rc[idx]
			char.CharType = "hancheng"
			e = fillRadicalDetail(radical, char)
			if e != nil {
				log.With("bushou", rc[idx].BuShou, "pinyin", rc[idx].PinYin).Error(e)
				continue
			}
		}
	case RadicalTypeHanChengBushou:
		rc = analyzeBushouRadical(document)
		bytes, e := json.Marshal(rc)
		if e != nil {
			return
		}
		log.With("size", len(rc)).Info(string(bytes))
		for idx := range rc {
			radical, e := RadicalReader(s, rc[idx].BuShou)
			if e != nil {
				return
			}
			char := rc[idx]
			char.CharType = "hancheng"
			e = fillRadicalDetail(radical, char)
			if e != nil {
				log.With("bushou", rc[idx].BuShou, "pinyin", rc[idx].PinYin).Error(e)
				continue
			}
		}
	}
}

func parseRadicalWD(radicalType RadicalType, character *RadicalCharacter, wd string) {
	switch radicalType {
	case RadicalTypeHanChengBushou, RadicalTypeKangXiBushou:
		//character.BuShou = wd
	case RadicalTypeHanChengPinyin, RadicalTypeKangXiPinyin:
		//character.Alphabet = wd
	case RadicalTypeHanChengBihua, RadicalTypeKangXiBihua:
		//character.BiHua = wd
	}
}

func fillRadicalDetail(radical *Radical, character *RadicalCharacter) (err error) {
	log.Infof("%+v", radical)
	for _, tmp := range *(*[]RadicalUnion)(radical) {
		for i := range tmp.RadicalCharacterArray {
			rc := tmp.RadicalCharacterArray[i]
			rc.Alphabet = character.Alphabet
			rc.BuShou = character.BuShou
			rc.CharType = character.CharType
			one, e := insertOrUpdateRadicalCharacter(db, &rc)
			if e != nil {
				return e
			}
			log.With("num", one).Info(rc)
		}
		log.With("value", radical).Info("radical")
	}

	return nil
}
