package excavator

import (
	"fmt"
	"log"
	"sync"
)

type CharacterFunc func(character *Character) error
type RadicalFunc func(radical *Radical) error

//RootRadical result root list
type Root struct {
	iterator
	wg     sync.WaitGroup
	URL    string
	Suffix string
}

//Radical
type Radical struct {
	root *Root
	iterator
	RadicalCharacter
}

type RadicalCharacter struct {
	Strokes string
	Pinyin  string
	Name    string
	URL     string
}


type Character struct {
	URL            string //汉字地址
	Character      string //字符
	Pinyin         string //拼音
	Radical        string //部首
	RadicalStrokes string //部首笔画
	TotalStrokes   string //总笔画
	KangxiStrokes  string //康熙笔画数
	Phonetic       string //注音
	Folk
	Structure
	Explain
	Rhyme
	Index
}

//民俗参考
type Folk struct {
	CommonlyCharacters   string //是否为常用字
	NameScience          string //姓名学
	FiveElementCharacter string //汉字五行
	GodBadMoral          string //吉凶寓意
}

//字形结构
type Structure struct {
	DecompositionSearch string //首尾分解查字
	StrokeNumber        string //笔顺编号
	StrokeReadWrite     string //笔顺读写
}

//Explain 康熙字典解释
type Explain struct {
	Intro  string //简介
	Detail string //详情
}

//Rhyme 音韵参考
type Rhyme struct {
	GuangYun  string //广　韵
	Mandarin  string //国　语
	Cantonese string //粤　语
}

//Index 索引参考
type Index struct {
	AncientWrite      string //古文字诂林
	HometownTrain     string //故训彙纂
	Explain           string //说文解字
	KangxiDictionary  string //康熙字典
	ChineseDictionary string //汉语字典
	Cihai             string //辞　海  
}

//NewRoot create an root
func NewRoot(url string, suffix string) *Root {
	return &Root{
		URL:    url,
		Suffix: suffix,
	}
}

func (r *Root) WaitForDone() {
	r.wg.Wait()
}

//Add add radical
func (r *Root) Add(rd *Radical) {
	r.iterator.Add(rd)
}

func (r *Root) Iterator(f RadicalFunc) {
	r.Reset()
	for r.HasNext() {
		radical := r.Next().(*Radical)
		if f(getRedicalList(r, radical)) != nil {
			break
		}
	}
}

//func (r *Root) ThreadIterator(f RadicalFunc) {
//	r.Reset()
//	wg := sync.WaitGroup{}
//	for r.HasNext() {
//		radical := r.Next().(*Radical)
//		wg.Add(1)
//		go func(group *sync.WaitGroup) {
//			defer group.Done()
//			if err := f(getRedicalList(r, radical)); err != nil {
//				panic(err)
//			}
//		}(&wg)
//	}
//	wg.Wait()
//}

func (root *Root) Self() *Root {
	return getRootList(root, root.Suffix)
}

func (root *Root) IteratorSelf(f RadicalFunc) {
	log.Println(root.iterator.Size())
	root.iterator.Iterator(func(v interface{}) error {
		v1 := v.(*Radical)
		return f(v1)
	})
}

func (root *Root) SelfRadical(name string) *Radical {
	var r *Radical
	root.Iterator(func(radical *Radical) error {
		if radical.Name == name {
			r = radical
			return fmt.Errorf("%s%s", " radical found:", radical.Name)
		}
		return nil
	})
	return r
}

func (root *Root) SelfRadicals() []*Radical {
	var radicals []*Radical
	root.Iterator(func(radical *Radical) error {
		radicals = append(radicals, radical)
		return nil
	})
	return radicals
}

//GetList get an list from web
func (root *Root) GetList() []*Character {
	var cs []*Character
	root.Self()
	rs := root.SelfRadicals()
	for _, r := range rs {
		cs = append(cs, r.SelfCharacters()...)
	}
	return cs
}

func (root *Root) ListProcess(f func(c *Character)) []*Character {
	var cs []*Character
	root.Self()
	rs := root.SelfRadicals()
	for _, r := range rs {
		for _, v := range r.SelfCharacters() {
			f(v)
		}
	}
	return cs
}

func (r *Radical) Add(character *Character) {
	r.iterator.Add(character)
}

func (r *Radical) Iterator(f CharacterFunc) {
	r.Reset()
	for r.HasNext() {
		character := r.Next().(*Character)
		if f(getCharacterList(r.root, character)) != nil {
			break
		}
	}
}

func (r *Radical) ThreadIterator(f CharacterFunc) {
	wg := sync.WaitGroup{}
	for r.HasNext() {
		character := r.Next().(*Character)
		wg.Add(1)
		go func(group *sync.WaitGroup) {
			defer group.Done()
			if err := f(getCharacterList(r.root, character)); err != nil {
				panic(err)
			}
		}(&wg)
	}
	wg.Wait()
}

func (r *Radical) SetRoot(root *Root) {
	r.root = root
}

func (r *Radical) SelfCharacter(name string) *Character {
	var c *Character
	r.Iterator(func(character *Character) error {
		if character.Character == name {
			c = character
			return fmt.Errorf("%s", "radical found")
		}
		return nil
	})
	return c
}

func (r *Radical) SelfCharacters() []*Character {
	var characters []*Character
	r.Iterator(func(character *Character) error {
		characters = append(characters, getCharacterList(r.root, character))
		return nil
	})
	return characters
}
