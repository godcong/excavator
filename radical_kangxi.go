package excavator

import (
	"fmt"
	"github.com/gocolly/colly"
)



func grabRadical(url string, characters chan<- *RadicalCharacter) {
	if characters == nil {
		return
	}
	defer func() {
		characters <- nil
	}()
	c := colly.NewCollector()
	c.OnHTML("a[href][data-action]", func(element *colly.HTMLElement) {
		da := element.Attr("data-action")
		log.With("value", da).Info("data action")
		if da == "" {
			return
		}
		q := NewQuery()

		closer, e := q.Grab(da)
		if e != nil {
			return
		}
		radical, e := RadicalReader(closer)
		if e != nil {
			return
		}
		for _, tmp := range *(*[]RadicalUnion)(radical) {
			for i := range tmp.RadicalCharacterArray {
				rc := tmp.RadicalCharacterArray[i]
				//e := exc.saveRadicalCharacter(&tmp.RadicalCharacterArray[i])
				//if e != nil {
				//	log.Error(e)
				//	continue
				//}
				characters <- &rc
			}
		}
		log.With("value", radical).Info("radical")
	})
	c.OnResponse(func(response *colly.Response) {
		log.Info(string(response.Body))
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})
	e := c.Visit(url)
	if e != nil {
		log.Error(e)
	}
	return
}
