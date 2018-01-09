package excavator

import (
	"bytes"
	"log"
	"net/http"
	"strings"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

var trans = false

func TransformOn() {
	trans = true
}

func TransformOff() {
	trans = false
}

func GetBodyFromUrl(url string) string {
	client := &http.Client{}
	req, e1 := http.NewRequest("GET", url, nil)
	if e1 != nil {
		return ""
	}

	resp, e2 := client.Do(req)
	if e2 != nil {
		return ""
	}

	defer resp.Body.Close()
	b := bytes.Buffer{}
	var e3 error
	if trans {
		reader := transform.NewReader(resp.Body, simplifiedchinese.GBK.NewDecoder())
		_, e3 = b.ReadFrom(reader)

	} else {
		_, e3 = b.ReadFrom(resp.Body)

	}

	if e3 != nil {
		return ""
	}

	return b.String()
}

func GetCharList(url string) map[int]map[string]string {
	mlist := make(map[int]map[string]string)
	s := GetBodyFromUrl(url)
	sa := StringSplite(s, "<tr bgcolor", "</tr>")

	for k, v := range sa {
		saa := DoFix(v)
		addrMap := DecodeChar(strings.NewReader(saa))
		mlist[k] = addrMap

	}

	return mlist
}

func GetRootList(url string) map[int]map[string]string {
	mlist := make(map[int]map[string]string)
	s := GetBodyFromUrl(url)
	sa := StringSplite(s, "<tr>", "</tr>")

	for k, v := range sa {
		saa := DoFix(v)
		log.Println(saa)
		addrMap := DecodeRoot(strings.NewReader(saa))
		mlist[k] = addrMap

	}

	return mlist
}

func SetFix(f func(string) string) {
	fix = f
}

func DoFix(s string) string {
	if fix != nil {
		return fix(s)
	}
	return s
}
