package net

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
)

// QueryGet ...
func QueryGet(url string) (*goquery.Document, error) {
	if cli == nil {
		cli = http.DefaultClient
	}

	res, err := cli.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	return doc, err
}

func CacheQuery(url string) (*goquery.Document, error) {
	reader, e := cache.Reader(url)
	if e != nil {
		return nil, e
	}
	//if cli == nil {
	//	cli = http.DefaultClient
	//}
	//
	//res, err := cli.Get(url)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer res.Body.Close()
	//if res.StatusCode != 200 {
	//	return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	//}

	return goquery.NewDocumentFromReader(reader)
}