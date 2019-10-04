package excavator

const HanChengBushou = "/bushou/zi/"
const HanChengPinyin = "/pinyin/zi/"
const HanChengBihua = "/bihua/zi/"

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
