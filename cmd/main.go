package main

import (
	"flag"

	"github.com/godcong/excavator"
)

var url = flag.String("url", "http://tool.httpcn.com", "catch the web url")

func main() {

	flag.Parse()

	local := *url
	suffix := "/KangXi/BuShou.html"
	m := local + suffix
	excavator.GetRootList(m)
}
