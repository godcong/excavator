package excavator

import (
	"fmt"
	"log"
)

type CharacterFunc func(character *Character) error
type RadicalFunc func(radical *Radical) error

//RootRadical result root list
type Root struct {
	iterator
	//Radicals []*Radical
	URL string
}

//Radical
type Radical struct {
	iterator
	Strokes string
	Name    string
	URL     string
	//Character []*Character
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

//音韵参考
type Explain struct {
}

//索引参考
type Rhyme struct {
	GuangYun  string //广　韵
	Mandarin  string //国　语
	Cantonese string //粤　语
}

type Index struct {
	AncientWrite      string //古文字诂林
	HometownTrain     string //故训彙纂
	Explain           string //说文解字
	KangxiDictionary  string //康熙字典
	ChineseDictionary string //汉语字典
	Cihai             string //辞　海  
}

var root *Root

func init() {
	root = NewRoot("http://tool.httpcn.com")
	root.Self("/KangXi/BuShou.html")
}

func Self() *Root {
	return root
}

func SelfRadical(name string) *Radical {
	return root.SelfRadical(name)
}

func SeflRadicals() []*Radical {
	return root.SelfRadicals()
}

//NewRoot create an root
func NewRoot(url string) *Root {
	return &Root{
		URL: url,
	}
}

//Add add radical
func (r *Root) Add(rd *Radical) {
	r.iterator.Add(rd)
}

func (r *Root) Iterator(f RadicalFunc) {
	r.Reset()
	for r.HasNext() {
		if f(r.Next().(*Radical)) != nil {
			break
		}
	}
}

func (root *Root) Self(s string) *Root {
	return getRootList(root, s)
}

func (root *Root) SelfRadical(name string) *Radical {
	var radical *Radical
	root.Iterator(func(radical *Radical) error {
		log.Println(radical.Name)
		if radical.Name == name {
			radical = getRedicalList(root, radical)
			return fmt.Errorf("%s%s", " radical found:", radical.Name)

		}
		return nil
	})
	return radical
}

func (root *Root) SelfRadicals() []*Radical {
	var radicals []*Radical
	root.Iterator(func(radical *Radical) error {
		radicals = append(radicals, getRedicalList(root, radical))
		return nil
	})
	return radicals
}

//GetList get an list from web
func (root *Root) GetList(s string) []*Character {
	var cs []*Character
	root.Self(s)
	rs := root.SelfRadicals()
	for _, r := range rs {
		cs = append(cs, r.SelfCharacters()...)
	}
	return cs
}

func (r *Radical) Add(character *Character) {
	r.iterator.Add(character)
}
func (r *Radical) Iterator(f CharacterFunc) {
	r.Reset()
	for r.HasNext() {
		if f(r.Next().(*Character)) != nil {
			break
		}
	}
}

func (r *Radical) SelfCharacter(name string) *Character {
	var c *Character
	r.Iterator(func(character *Character) error {
		if character.Character == name {
			c = getCharacterList(Self(), character)
			return fmt.Errorf("%s", "radical found")
		}
		return nil
	})
	return c
}

func (r *Radical) SelfCharacters() []*Character {
	var characters []*Character
	r.Iterator(func(character *Character) error {
		characters = append(characters, getCharacterList(Self(), character))
		return nil
	})
	return characters
}
