package main

import (
	"github.com/go-xorm/xorm"
	"github.com/godcong/excavator"
	"github.com/godcong/excavator/net"
	"github.com/godcong/fate"
	"github.com/godcong/fate/config"
)

func main() {
	fromDB := excavator.InitMysql("127.0.0.1:3306", "root", "111111")

	db := fate.InitDatabaseFromConfig(config.Config{
		Database: config.Database{
			Host:         "localhost",
			Port:         "3306",
			User:         "root",
			Pwd:          "111111",
			Name:         "fate",
			MaxIdleCon:   0,
			MaxOpenCon:   0,
			Driver:       "mysql",
			File:         "",
			Dsn:          "",
			ShowSQL:      true,
			ShowExecTime: false,
		},
	})

	chars := make(chan *excavator.Character)

	go getCharacters(fromDB, chars)

	for char := range chars {

		fc, e := db.GetCharacter(fate.Char(char.Ch))

		if e != nil {
			c := fate.Character{
				Hash:                     net.Hash(char.Ch),
				PinYin:                   char.PinYin,
				Ch:                       char.Ch,
				Radical:                  char.Radical,
				RadicalStroke:            char.RadicalStroke,
				Stroke:                   char.Stroke,
				IsKangXi:                 char.IsKangXi,
				KangXi:                   char.KangXi,
				KangXiStroke:             char.KangXiStroke,
				SimpleRadical:            char.SimpleRadical,
				SimpleRadicalStroke:      char.SimpleRadicalStroke,
				SimpleTotalStroke:        char.SimpleTotalStroke,
				TraditionalRadical:       char.TraditionalRadical,
				TraditionalRadicalStroke: char.TraditionalRadicalStroke,
				TraditionalTotalStroke:   char.TraditionalTotalStroke,
				NameScience:              char.NameScience,
				WuXing:                   char.WuXing,
				Lucky:                    char.Lucky,
				Regular:                  char.Regular,
				TraditionalCharacter:     char.TraditionalCharacter,
				VariantCharacter:         char.VariantCharacter,
				Comment:                  char.Comment,
			}
			fixStroke(&c)
			_, e := db.Database().(*xorm.Engine).InsertOne(&c)
			if e != nil {
				panic(e)
			}
			continue
		}
		fc.Hash = net.Hash(char.Ch)
		fc.PinYin = char.PinYin
		fc.Ch = char.Ch
		fc.Radical = char.Radical
		fc.RadicalStroke = char.RadicalStroke
		fc.Stroke = char.Stroke
		fc.IsKangXi = char.IsKangXi
		fc.KangXi = char.KangXi
		fc.KangXiStroke = char.KangXiStroke
		fc.SimpleRadical = char.SimpleRadical
		fc.SimpleRadicalStroke = char.SimpleRadicalStroke
		fc.SimpleTotalStroke = char.SimpleTotalStroke
		fc.TraditionalRadical = char.TraditionalRadical
		fc.TraditionalRadicalStroke = char.TraditionalRadicalStroke
		fc.TraditionalTotalStroke = char.TraditionalTotalStroke
		fc.NameScience = char.NameScience
		fc.WuXing = char.WuXing
		fc.Lucky = char.Lucky
		fc.Regular = char.Regular
		fc.TraditionalCharacter = char.TraditionalCharacter
		fc.VariantCharacter = char.VariantCharacter
		fc.Comment = char.Comment
		fixStroke(fc)
		_, e = db.Database().(*xorm.Engine).Where("hash =?", fc.Hash).Update(fc)
		if e != nil {
			panic(e)
		}
	}

}

func getCharacters(engine *xorm.Engine, c chan<- *excavator.Character) (e error) {
	rows, e := engine.Rows(&excavator.Character{})
	if e != nil {
		return e
	}

	for rows.Next() {
		var c excavator.Character
		e := rows.Scan(&c)
		if e != nil {
			return e
		}
	}
	close(c)
	return nil
}
func fixStroke(character *fate.Character) bool {
	if character.KangXiStroke != 0 {
		character.ScienceStroke = character.KangXiStroke
	} else if character.TraditionalTotalStroke != 0 {
		character.ScienceStroke = character.TraditionalTotalStroke
	} else if character.Stroke != 0 {
		character.ScienceStroke = character.Stroke
	} else if character.SimpleTotalStroke != 0 {
		character.ScienceStroke = character.SimpleTotalStroke
	} else {
		return false
	}
	return true
}
