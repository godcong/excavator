package excavator_test

import (
	"log"
	"testing"

	"github.com/godcong/excavator"
)

func TestGet(t *testing.T) {
	l := excavator.Get("http://xh.5156edu.com")
	log.Println(l)
}

func TestGetKangXi(t *testing.T) {
	excavator.GetKangXi("http://tool.httpcn.com")
}