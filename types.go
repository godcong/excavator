package excavator

import (
	"log"
)

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
		getRedicalList(root, radical)
	})
	root.Iterator(func(radical *Radical) {
		radical.Iterator(func(character *Character) {
			log.Println(*character)
			getCharacterList(root, character)
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
