package excavator

import (
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/mattn/go-sqlite3"
	"net/url"
)

var _ = &sqlite3.SQLiteDriver{}
var _ = &mysql.Config{}

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

const sqlURL = "%s:%s@tcp(%s)/%s?loc=%s&charset=utf8mb4&parseTime=true"

func InitMysql(addr, name, pass string) *xorm.Engine {
	u := fmt.Sprintf(sqlURL, name, pass, addr, "excavator", url.QueryEscape("Asia/Shanghai"))
	eng, e := xorm.NewEngine("mysql", u)
	if e != nil {
		panic(e)
	}
	eng.ShowSQL(true)
	eng.ShowExecTime(true)
}
