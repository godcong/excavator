package excavator_test

import (
	"log"
	"testing"

	"github.com/godcong/excavator"
)

func TestRadical_Iterator(t *testing.T) {
	root := excavator.Self()
	radical := root.SelfRadical("丨")
	log.Println(radical)
	c := radical.SelfCharacters()
	log.Println(c)
}
