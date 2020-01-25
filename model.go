package excavator

import (
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/mattn/go-sqlite3"
	"github.com/xormsharp/xorm"
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
		log.Fatal(e)
	}
	return eng
}

const sqlURL = "%s:%s@tcp(%s)/%s?loc=%s&charset=utf8mb4&parseTime=true"

func InitMysql(addr, name, pass string) *xorm.Engine {
	u := fmt.Sprintf(sqlURL, name, pass, addr, "excavator", url.QueryEscape("Asia/Shanghai"))
	eng, e := xorm.NewEngine("mysql", u)
	if e != nil {
		log.Fatal(e)
	}
	eng.ShowSQL(true)
	eng.ShowExecTime(true)
	return eng
}
