package db_test

import (
	"testing"

	"github.com/godcong/excavator"
	"github.com/godcong/excavator/db"
)

func TestInsertRootFromJson(t *testing.T) {
	db.InsertRootFromJson("root.json", "root")
}

func TestInsertRadicalFromJson(t *testing.T) {
	db.InsertRadicalFromJson("radical.json", "radical2")
}

func TestInsertFromJson(t *testing.T) {
	db.InsertFromJson("wuxing.json", &db.WuXing{})
}

func TestInsertWuXingFromJson(t *testing.T) {
	db.InsertWuXingFromJson("wuxing.json", "wuxing")
}

func TestInsertDaYanFromJson(t *testing.T) {
	db.InsertDaYanFromJson("dayan.json", "dayan")
}

func TestInsertCharacterFromJson(t *testing.T) {
	db.InsertCharacterFromJson("character.json", "character")
}

func TestUpdateCommonly(t *testing.T) {
	//chars := excavator.CommonlyTop("http://www.zdic.net/z/zb/cc1.htm")
	chars := excavator.CommonlyTop("http://www.zdic.net/z/zb/cc2.htm")


	for _, v := range chars {
		//bc := excavator.CommonlyBase("http://www.zdic.net", v)
		db.DB("commonly").Insert(v)
	}

}
