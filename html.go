package excavator

import (
	"bytes"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

type CharDetail struct {
	Char           string
	NameType       string
	NameRoot       string
	Pinyin         string
	Radical        string
	SimpleStrokes  int
	ScienceStrokes int
}

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
	s = strings.ToLower(s)
	sa := StringSplite(s, "<tr bgcolor=#ffffff", "</tr>", true)

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
	sa := StringSplite(s, "<tr>", "</tr>", true)

	for k, v := range sa {
		saa := DoFix(v)
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

func GetFileterCharList(url string) map[int]map[string]string {
	mlist := make(map[int]map[string]string)
	s := GetBodyFromUrl(url)
	s = strings.ToLower(s)
	sa := StringSplite(s, `<tr bgcolor=#ffffff>`, "</tr>", true)

	for k, v := range sa {
		saa := DoFix(v)
		addrMap := DecodeRoot(strings.NewReader(saa))
		mlist[k] = addrMap

	}

	return mlist
}

func GetCharDetail(url string) CharDetail {
	var c, t, r string
	var b int
	s := GetBodyFromUrl(url)
	s = strings.ToLower(s)

	hz := StringSplite(s, `『`, "』", false)
	if len(hz) == 0 {
		return CharDetail{}
	}
	if len(hz) > 0 {
		c = hz[0]
	}
	wx := StringSplite(s, `汉字五行：`, "　", false)
	if len(wx) > 0 {
		t = wx[0]
	}
	fj := StringSplite(s, `首尾分解查字</span> ]：`, "(", false)
	if len(fj) > 0 {
		r = fj[0]
	}
	bh := StringSplite(s, `康熙笔画：`, "；", false)
	if len(bh) > 0 {
		b, _ = strconv.Atoi(bh[0])
	}
	return CharDetail{
		Char: c, NameType: t, NameRoot: r, ScienceStrokes: b,
	}
}
