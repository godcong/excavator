package excavator_test

import (
	"log"
	"testing"

	"github.com/godcong/excavator"
)

func TestIterator_Add(t *testing.T) {
	iter := excavator.NewIterator()
	iter.Add(8)
	log.Println(iter)
}
