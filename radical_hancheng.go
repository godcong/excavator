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
	//radical := make(chan *RadicalCharacter)
	rc := analyzeRadical(document)
	log.With("size", len(rc)).Info("radicals")

	for _, r := range rc {
		//r := <-radical
		//if r == nil {
		//	break
		//}
		log.Info(fillRadicalDetail(r))
	}
}

func analyzeRadical(document *goquery.Document) (rc []*RadicalCharacter) {
	//rc := make(chan *RadicalCharacter)
	//defer func() {
	//	rc <- nil
	//}()

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
		rc = append(rc, ch)
	})
	return
}

func fillRadicalDetail(character *RadicalCharacter) (*RadicalCharacter, error) {
	q := NewQuery(RequestTypeOption(RequestTypeHanCheng))

	closer, e := q.Grab(character.BuShou)
	if e != nil {
		return nil, e
	}

	radical, e := RadicalReader(closer)
	if e != nil {
		return nil, e
	}
	log.Infof("%+v", radical)
	for _, tmp := range *(*[]RadicalUnion)(radical) {
		for i := range tmp.RadicalCharacterArray {
			rc := tmp.RadicalCharacterArray[i]
			one, e := insertRadicalCharacter(db, &tmp.RadicalCharacterArray[i])
			if e != nil {
				log.Error(e)
				continue
			}
			log.With("num", one).Info(rc)
		}
		log.With("value", radical).Info("radical")
	}

	return new(RadicalCharacter), nil
}
