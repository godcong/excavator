package excavator

import (
	"encoding/json"
	"github.com/godcong/excavator/net"
)

const HanChengBushou = "/bushou/zi/"
const HanChengPinyin = "/pinyin/zi/"
const HanChengBihua = "/bihua/zi/"

func grabHanChengRadicalList(s RadicalType, url string) {
	document, e := net.CacheQuery(url)
	if e != nil {
		panic(e)
	}

	switch s {
	case RadicalTypeHanChengPinyin:
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
	case RadicalTypeHanChengBushou:
		rc := analyzeBushouRadical(document)
		bytes, e := json.Marshal(rc)
		if e != nil {
			return
		}
		log.With("size", len(rc)).Info(string(bytes))

		for idx := range rc {
			e := fillRadicalBushouDetail(rc[idx])
			if e != nil {
				log.With("bushou", rc[idx].BuShou, "pinyin", rc[idx].PinYin).Error(e)
				continue
			}
		}
		//case SearchTypeBihua:

	}
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
func fillRadicalBushouDetail(character *RadicalCharacter) (err error) {
	q := NewQuery(RequestTypeOption(RequestTypeHanChengBushou))

	closer, e := q.Grab(character.BuShou)
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
