package excavator

import (
	"encoding/xml"
	"io"
	"log"
)

var fix func(string) string

func DecodeRoot(r io.Reader) map[string]string {
	var t xml.Token
	var err error
	//var sa []string
	addrMap := make(map[string]string)
	addA := false
	tmp := ""

	d := xml.NewDecoder(r)
	t, err = d.Token()
	for ; err == nil; t, err = d.Token() {

		switch token := t.(type) {
		// 处理元素开始（标签）
		case xml.StartElement:
			name := token.Name.Local
			if name == "a" {
				addA = true
			}

			for _, attr := range token.Attr {
				attrName := attr.Name.Local
				attrValue := attr.Value
				if attrName == "href" {
					tmp = attrValue
				}

			}
			// 处理元素结束（标签）
		case xml.EndElement:
			name := token.Name.Local
			if name == "a" {
				//addrMap[tmp] = sa
				tmp = ""
				//sa = []string{}
				addA = false
			}
			// 处理字符数据（这里就是元素的文本）

		case xml.CharData:

			content := string([]byte(token))
			if addA {
				//sa = append(sa, content)
				addrMap[tmp] = content
				tmp = ""
			}

		default:

			// ...

		}

	}
	if err != nil && err != io.EOF {
		log.Println("error", err)
	}
	return addrMap
}

func DecodeChar(r io.Reader) map[string]string {
	var t xml.Token
	var err error

	addrMap := make(map[string]string)
	addA := false
	addSpan := false
	tmp := ""
	zf := ""
	py := ""
	bh := ""
	d := xml.NewDecoder(r)
	t, err = d.Token()
	countTD := 0

	for ; err == nil; t, err = d.Token() {

		switch token := t.(type) {
		// 处理元素开始（标签）
		case xml.StartElement:

			name := token.Name.Local

			//fmt.Printf("Token name: %s\n", name)
			if name == "a" {
				addA = true
			}

			if name == "span" {
				addSpan = true
			}
			for _, attr := range token.Attr {

				attrName := attr.Name.Local

				attrValue := attr.Value
				//fmt.Printf("An attribute is: %s %s\n", attrName, attrValue)
				if attrName == "href" {
					tmp = attrValue
				}

			}

			// 处理元素结束（标签）

		case xml.EndElement:
			name := token.Name.Local
			//fmt.Printf("Token of '%s' end\n", name)
			if name == "td" {
				countTD++
			}

			if name == "a" {
				addA = false
			}

			if name == "span" {
				addSpan = false
			}
			// 处理字符数据（这里就是元素的文本）

		case xml.CharData:
			content := string([]byte(token))
			if countTD == 0 {
				bh = content
			}

			if addSpan {
				py = content
				continue
			}
			//if addP {
			//	bh = content
			//	continue
			//}

			if addA {
				zf = content
				continue
			}

			//fmt.Printf("This is the content: %v\n", content)

		default:
			// ...
		}
		if bh != "" && zf != "" && py != "" {
			addrMap[tmp] = bh + "|" + zf + "|" + py
			zf = ""
			py = ""
		}
	}

	return addrMap
}
