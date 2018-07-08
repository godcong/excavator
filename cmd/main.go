package main

import (
	"flag"
	"log"

	"github.com/godcong/excavator"
)

var url = flag.String("url", "http://tool.httpcn.com", "catch the web url")

func main() {

	flag.Parse()

	local := *url
	suffix := "/KangXi/BuShou.html"
	log.SetFlags(log.Llongfile)
	root := excavator.NewRoot(local)
	root.GetList(suffix)
}
