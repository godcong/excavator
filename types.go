package excavator

import (
	"log"
)

type iterator struct {
	data  []interface{}
	index int
}

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
	KangxiStrokes  string //康熙笔画数
	Phonetic       string //注音
	Folk
}

type Folk struct {
	CommonlyCharacters   string //是否为常用字
	NameScience          string //姓名学
	FiveElementCharacter string //汉字五行
	GodBadMoral          string //吉凶寓意
}

type Structure struct {
}

type Explain struct {
}

type Rhyme struct {
}

type Index struct {
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

func (r *Root) Iterator(f func(radical *Radical)) {
	r.Reset()
	if r.HasNext() {
		f(r.Next().(*Radical))
	}
}

//GetList get an list from web
func (root *Root) GetList(s string) {
	getRootList(root, s)
	//wg := sync.WaitGroup{}
	root.Iterator(func(radical *Radical) {
		log.Println(radical)
		getRedicalList(root,radical)
	})
	root.Iterator(func(radical *Radical) {
		radical.Iterator(func(character *Character) {
			getCharacterList(root,character)
		})

	})
}

func (r *Radical) Add(character *Character) {
	r.iterator.Add(character)
}
func (r *Radical) Iterator(f func(character *Character)) {
	r.Reset()
	if r.HasNext() {
		f(r.Next().(*Character))
	}
}

//HasNext check next
func (i *iterator) HasNext() bool {
	return i.index < len(i.data)
}

//Next get next
func (i *iterator) Next() interface{} {
	defer func() {
		i.index++
	}()
	if i.index < len(i.data) {
		return i.data[i.index]
	}

	return nil
}

//Reset reset index
func (i *iterator) Reset() {
	i.index = 0
}

//Add add radical
func (i *iterator) Add(v interface{}) {
	i.data = append(i.data, v)
}
