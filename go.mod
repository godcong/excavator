module github.com/godcong/excavator

require (
	github.com/antchfx/htmlquery v1.2.5-0.20211125074323-810ee8082758
	github.com/free-utils-go/xorm_type_assist v0.0.0-20210507214645-5c65059bdc6a
	github.com/gin-gonic/gin v1.7.2-0.20220214063957-375714258462
	github.com/go-sql-driver/mysql v1.6.1-0.20220303001332-c1aa6812e475
	github.com/godcong/cachenet v0.0.0-00010101000000-000000000000
	github.com/godcong/fate v0.0.0-00010101000000-000000000000
	github.com/godcong/go-trait v0.0.0-20190816072228-f216e906756e
	github.com/goextension/log v0.0.2
	github.com/mattn/go-sqlite3 v1.14.12
	github.com/pkg/browser v0.0.0-20210911075715-681adbf594b8
	github.com/tebeka/selenium v0.9.10-0.20211105214847-e9100b7f5ac1
	golang.org/x/exp v0.0.0-20220307200941-a1099baf94bf
	golang.org/x/net v0.0.0-20220225172249-27dd8689420f
	xorm.io/xorm v1.2.4-0.20220125052846-3180c418c245
)

go 1.16

// pseudo-version can be got from `go get  github.com/<name>/<project>@<commit>`
// When merged into upstream, the replace sentence can be disabled

replace (
	github.com/godcong/cachenet => github.com/free-utils-go/cachenet v0.0.0-20210507230213-b162ba0f9c57
	github.com/godcong/fate => github.com/fortune-fun/fate v0.0.0-20210519004006-1dbc3506aef8
	github.com/godcong/name_gua => github.com/fortune-fun/name_gua v0.0.0-20210515180506-8c0f084200f1
	github.com/godcong/name_wuge => github.com/fortune-fun/name_wuge v0.0.0-20210510141111-8cee898249c6
	github.com/godcong/yi => github.com/fortune-fun/yi v0.0.0-20210518235908-d42db1a65871
)
