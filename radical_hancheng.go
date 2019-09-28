package excavator

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/godcong/excavator/net"
)

const HanChengMainPage = "http://hy.httpcn.com/bushou/zi/"

func grabRadicalList(url string) {
	document, e := net.CacheQuery(url)
	if e != nil {
		return
	}
	log.Info(document.Text())
}

func analyzeRadical(document *goquery.Document) chan<- *RadicalCharacter {
	rc := make(chan *RadicalCharacter)
	document.Find("")
}
