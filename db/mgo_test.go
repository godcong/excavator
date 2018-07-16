package db_test

import (
	"testing"

	"github.com/godcong/excavator/db"
)

func TestInsertRootFromJson(t *testing.T) {
	db.InsertRootFromJson("root.json")
}

func TestInsertRadicalFromJson(t *testing.T) {
	db.InsertRadicalFromJson("radical.json")
}