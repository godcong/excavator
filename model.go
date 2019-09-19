package excavator

import (
	"github.com/go-xorm/xorm"
	"github.com/mattn/go-sqlite3"
)

var _ = &sqlite3.SQLiteDriver{}

// InitSqlite3 ...
func InitSqlite3(name string) *xorm.Engine {
	eng, e := xorm.NewEngine("sqlite3", name)
	if e != nil {
		panic(e)
	}
	eng.ShowSQL(true)
	eng.ShowExecTime(true)
	_, e = eng.Exec("PRAGMA journal_mode = OFF;")
	if e != nil {
		return nil
	}
	return eng
	//log.Info("result:", result)
	//for idx, val := range syncTable {
	//	log.Info("syncing ", idx)
	//	e := eng.Sync2(val)
	//	if e != nil {
	//		return e
	//	}
	//}
	//
	//db = eng
	//return nil
}
