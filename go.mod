module github.com/godcong/excavator

require (
	github.com/antchfx/htmlquery v1.2.3
	github.com/free-utils-go/cachenet v0.0.0-20210507230213-b162ba0f9c57
	github.com/free-utils-go/xorm_type_assist v0.0.0-20210507214645-5c65059bdc6a
	github.com/gin-gonic/gin v1.5.0
	github.com/go-sql-driver/mysql v1.5.0
	github.com/godcong/fate v0.0.0-00010101000000-000000000000
	github.com/godcong/go-trait v0.0.0-20190816072228-f216e906756e
	github.com/goextension/log v0.0.2
	github.com/mattn/go-sqlite3 v1.14.0
	github.com/pkg/browser v0.0.0-20180916011732-0a3d74bf9ce4
	github.com/tebeka/selenium v0.9.4
	golang.org/x/exp v0.0.0-20191030013958-a1ab85dbe136
	golang.org/x/net v0.0.0-20210428140749-89ef3d95e781
	xorm.io/xorm v1.0.7
)

go 1.16

replace (
	github.com/godcong/fate => ../fate
	github.com/godcong/yi => ../yi
)
