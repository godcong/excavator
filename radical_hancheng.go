package excavator

import (
	"encoding/json"
	"github.com/godcong/excavator/net"
)

const HanChengBushou = "/bushou/zi/"
const HanChengPinyin = "/pinyin/zi/"
const HanChengBihua = "/bihua/zi/"

func grabRadicalList(s RadicalType, url string) (e error) {
	document, e := net.CacheQuery(getMainURL(s, url))
	if e != nil {
		return e
	}
	var rc []*RadicalCharacter
	switch s {
	case RadicalTypeHanChengPinyin:
		rc = analyzePinyinRadical(document)
		bytes, e := json.Marshal(rc)
		if e != nil {
			return e
		}
		log.With("size", len(rc)).Info(string(bytes))
		for idx := range rc {
			radical, e := RadicalReader(s, rc[idx].PinYin, "")
			if e != nil {
				return e
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
			return e
		}
		log.With("size", len(rc)).Info(string(bytes))
		for idx := range rc {
			radical, e := RadicalReader(s, rc[idx].BuShou, "")
			if e != nil {
				return e
			}
			char := rc[idx]
			char.CharType = "hancheng"
			e = fillRadicalDetail(radical, char)
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
			radical, e := RadicalReader(s, rc[idx].BHNum, rc[idx].QBNum)
			if e != nil {
				return e
			}
			char := rc[idx]
			char.CharType = "hancheng"
			e = fillRadicalDetail(radical, char)
			if e != nil {
				log.With("bushou", rc[idx].BuShou, "pinyin", rc[idx].PinYin, "bihua", rc[idx].BiHua).Error(e)
				continue
			}
		}
	}
	return nil
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
			rc.BiHua = character.BiHua
			rc.QiBi = character.QiBi
			rc.QBNum = character.QBNum
			rc.BHNum = character.BHNum
			rc.TotalBiHua = character.TotalBiHua
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
