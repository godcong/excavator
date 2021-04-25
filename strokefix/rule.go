package strokefix

import (
	"excavator"

	"github.com/godcong/fate"
)

var numberCharList = `一二三四五六七八九十`

func NumberChar(exc *excavator.Excavator) {
	numberCharList_ := []rune(numberCharList)
	for idx, num_char := range numberCharList_ {
		my_char := fate.Character{
			Ch: string(num_char),
		}

		exc.DbFate.Get(&my_char)

		my_char.ScienceStroke = 1 + idx

		excavator.InsertOrUpdate(exc.DbFate, &my_char)
	}
}
