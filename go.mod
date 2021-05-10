module github.com/godcong/excavator

require (
	github.com/antchfx/htmlquery v1.2.3
	github.com/free-utils-go/xorm_type_assist v0.0.0-20210507214645-5c65059bdc6a
	github.com/gin-gonic/gin v1.5.0
	github.com/go-sql-driver/mysql v1.5.0
	github.com/godcong/cachenet v0.0.0-00010101000000-000000000000
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

// pseudo-version can be got from `go get  github.com/<name>/<project>@<commit>`

replace (
	github.com/godcong/cachenet => github.com/free-utils-go/cachenet v0.0.0-20210507230213-b162ba0f9c57
	github.com/godcong/fate => github.com/fortune-fun/fate v0.0.0-20210510144152-af89232cea41
	github.com/godcong/name_gua => github.com/fortune-fun/name_gua v0.0.0-20210510140743-082af5cba3cc
	github.com/godcong/name_wuge => github.com/fortune-fun/name_wuge v0.0.0-20210510141111-8cee898249c6
	github.com/godcong/yi => github.com/fortune-fun/yi v0.0.0-20210510135217-e2095161b447
)
