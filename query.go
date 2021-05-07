package excavator

import (
	"bufio"
	"errors"
	"excavator/models"
	"fmt"
	"math/bits"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/antchfx/htmlquery"
	"github.com/godcong/cachenet"
	"golang.org/x/net/html"
)

//从汉程获取文字的URL
func findHanChengUrl(uc *models.UnihanChar, html_node *html.Node) (string, error) {

	char_a := htmlquery.FindOne(html_node, "//ul[contains(@class, 'sotab_zi')]/li/a/span[normalize-space(text())='"+uc.UnicodeHex+"']/..")

	if char_a == nil {
		return "", errors.New("未找到url")
	}

	url := htmlquery.SelectAttr(char_a, "href")

	return url, nil
}

//从Unicode列表获取数据
func getHanChengChars(exc *Excavator, c chan<- *models.UnihanChar) (e error) {
	f, err := os.Open(exc.unicode_file)

	if err != nil {
		Log.Fatal(err)
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()
		unicode_hex := strings.Split(line, "U+")[1]
		unicode_int64, err := strconv.ParseInt(unicode_hex, 16, bits.UintSize)
		if err != nil {
			panic(line)
		}

		uc := models.UnihanChar{
			Unid:       int(unicode_int64),
			UnicodeHex: unicode_hex,
		}

		c <- &uc
	}

	if err := scanner.Err(); err != nil {
		Log.Fatal(err)
	}

	close(c)
	return nil
}

//从Unicode列表数据去汉程查询
func grabHanChengList(exc *Excavator) (err error) {
	chars := make(chan *models.UnihanChar)

	go getHanChengChars(exc, chars)

	for char := range chars {

		err = InsertOrUpdate(exc.Db, char)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}

		data := url.Values{
			"Tid": {"10"},
			"wd":  {char.UnicodeHex},
		}
		dataStr := data.Encode()

		html_node, err := cachenet.CacheDataQuery(exc.base_url, dataStr)
		if err != nil {
			panic(err)
		}

		url, err := findHanChengUrl(char, html_node)
		if err != nil {
			fmt.Println(err)
			continue
		}

		hcc := models.HanChengChar{
			Unid: char.Unid,
			Url:  url,
		}

		fmt.Println(url)
		//DB insert
		err = InsertOrUpdate(exc.Db, &hcc)
		if err != nil {
			panic(err)
		}
	}

	return nil
}
