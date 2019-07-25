package excavator

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/PuerkitoBio/goquery"
	"github.com/godcong/go-trait"
)

var log = trait.NewZapSugar()

const tmpFile = "tmp"

// Excavator ...
type Excavator struct {
	Workspace string `json:"workspace"`
	URL       string `json:"url"`
	HTML      string `json:"html"`
	Radical   map[string]string
}

// New ...
func New(url string, workspace string) *Excavator {
	return &Excavator{URL: url, Workspace: workspace}
}

// Run ...
func (exc *Excavator) Run() error {
	return exc.parseRadical()
}

func (exc *Excavator) parseRadical() (e error) {
	doc, e := exc.parseDocument(exc.URL)
	if e != nil {
		return e
	}
	exc.HTML, e = doc.Html()
	if e != nil {
		return e
	}
	return nil
}

//ParseDocument get the url result body
func (exc *Excavator) parseDocument(url string) (doc *goquery.Document, e error) {
	var reader io.Reader
	hash := SHA256(url)
	log.Infof("hash:%s,url:%s", hash, url)
	if !exc.IsExist(hash) {
		// Request the HTML page.
		res, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()
		if res.StatusCode != 200 {
			return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
		}
		reader = res.Body
		file, e := os.OpenFile(exc.getFilePath(hash), os.O_RDWR|os.O_CREATE|os.O_SYNC, os.ModePerm)
		if e != nil {
			return nil, e

		}
		written, e := io.Copy(file, reader)
		if e != nil {
			return nil, e
		}
		log.Infof("read %s | %d ", hash, written)
		_ = file.Close()
	}
	reader, e = os.Open(exc.getFilePath(hash))
	if e != nil {
		return nil, e
	}
	// Load the HTML document
	return goquery.NewDocumentFromReader(reader)
}

// IsExist ...
func (exc *Excavator) IsExist(name string) bool {
	_, e := os.Open(name)
	return e == nil || os.IsExist(e)
}

// GetPath ...
func (exc *Excavator) getFilePath(s string) string {
	if exc.Workspace == "" {
		exc.Workspace, _ = os.Getwd()
	}
	log.With("workspace", exc.Workspace, "tmpFile", tmpFile, s)
	return filepath.Join(exc.Workspace, tmpFile, s)
}

// SHA256 ...
func SHA256(s string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(s)))
}
