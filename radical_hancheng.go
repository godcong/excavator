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
