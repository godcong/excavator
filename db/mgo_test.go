package db_test

import (
	"testing"

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
