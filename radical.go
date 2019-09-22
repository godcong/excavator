package excavator

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"strings"
)

// Radical ...
type Radical []RadicalUnion

// UnmarshalRadical ...
func UnmarshalRadical(data []byte) (*Radical, error) {
	var r Radical
	err := json.Unmarshal(data, &r)
	return &r, err
}

// Marshal ...
func (r *Radical) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// RadicalCharacter ...
type RadicalCharacter struct {
	Zi     string `json:"zi"`
	PinYin string `json:"pinyin"`
	BuShou string `json:"bushou"`
	Num    string `json:"num"`
	URL    string `json:"url" xorm:"url"`
}

// RadicalUnion ...
type RadicalUnion struct {
	String                *string
	RadicalCharacterArray []RadicalCharacter
}

// UnmarshalJSON ...
func (x *RadicalUnion) UnmarshalJSON(data []byte) error {
	x.RadicalCharacterArray = nil
	_, err := unmarshalUnion(data, nil, nil, nil, &x.String, true, &x.RadicalCharacterArray, false, nil, false, nil, false, nil, false)
	if err != nil {
		return err
	}
	return nil
}

// MarshalJSON ...
func (x *RadicalUnion) MarshalJSON() ([]byte, error) {
	return marshalUnion(nil, nil, nil, x.String, x.RadicalCharacterArray != nil, x.RadicalCharacterArray, false, nil, false, nil, false, nil, false)
}

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
		q := NewQuery(url)

		r, e := q.AJAX(strings.NewReader(fmt.Sprintf("wd=%s", da)))
		if e != nil {
			return
		}
		radical, e := UnmarshalRadical(r)
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
		log.With("value", r).Info("radical")
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
