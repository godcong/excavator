package excavator

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/godcong/excavator/net"
)

const HanChengMainPage = "http://hy.httpcn.com/bushou/zi/"

func grabRadicalList(url string) {
	document, e := net.CacheQuery(url)
	if e != nil {
		panic(e)
	}
	analyzeRadical(document)
}

func analyzeRadical(document *goquery.Document) chan<- *RadicalCharacter {
	rc := make(chan *RadicalCharacter)
	document.Find("#segmentedControls > ul > li.mui-table-view-cell.mui-collapse").Each(func(i int, selection *goquery.Selection) {
		log.Info(selection.Html())
		ch := new(RadicalCharacter)
		ch.BiHua = selection.Find("a.mui-navigate-right").Text()
		selection.Find("div > a[data-action]").Each(func(i int, selection *goquery.Selection) {
			log.With("index", i, "text", selection.Text()).Info("bushou")
			bushouChar := *ch
			bushouChar.BuShou, _ = selection.Attr("data-action")
			log.With("bushou", bushouChar.BuShou).Info("bushou")
		})
		log.Infof("radical[%+v]", *ch)
	})
	return rc
}

func fillRadicalDetail(character *RadicalCharacter) *RadicalCharacter {
	return new(RadicalCharacter)
}
