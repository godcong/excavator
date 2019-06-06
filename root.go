package excavator

import (
	"github.com/godcong/go-trait"
	"go.uber.org/zap"
)

var log = trait.NewZapSugar()

// Debug ...
func Debug() {
	log = zap.NewExample(zap.AddCaller()).Sugar()
}

// RootFunc ...
type RootFunc func(rc *RootCharacter) error

//Root result root list
type Root struct {
	beforeIterator RootFunc
	iterator
	URL    string
	Suffix string
}

// RootCharacter ...
type RootCharacter struct {
	Strokes string `json:"strokes"`
	Name    string `json:"name"`
	URL     string `json:"url"`
}

//NewRoot create an root
func NewRoot(url string, suffix string) *Root {
	return &Root{
		URL:    url,
		Suffix: suffix,
	}
}

//Add add radical
func (r *Root) Add(rd *RootCharacter) {
	r.iterator.Add(rd)
}

// Self ...
func (r *Root) Self() *Root {
	return getRootList(r, r.Suffix)
}

// IteratorFunc ...
func (r *Root) IteratorFunc(f RootFunc) []*Radical {
	var rad []*Radical
	r.Reset()
	for r.HasNext() {
		rc := r.Next().(*RootCharacter)
		if r.beforeIterator != nil {
			if err := r.beforeIterator(rc); err != nil {
				log.Panic(err)
				continue
			}
		}
		rad = append(rad, getRedicalList(r, rc))
		if err := f(rc); err != nil {
			log.Panic(err)
			continue
		}
	}
	return rad
}

// Radical ...
func (r *Root) Radical(character *RootCharacter) *Radical {
	return getRedicalList(r, character)
}

// SetBefore ...
func (r *Root) SetBefore(rf RootFunc) {
	r.beforeIterator = rf
}
