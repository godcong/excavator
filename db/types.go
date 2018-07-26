package db

import "gopkg.in/mgo.v2/bson"

type IteratorFunc func(v interface{}) error

type WuXing struct {
	ID      bson.ObjectId `bson:"_id,omitempty"`
	WuXing  []string      `bson:"wu_xing" json:"wu_xing"`
	Fortune string        `bson:"fortune" json:"fortune"`
	Comment string        `bson:"comment" json:"comment"`
}

type DaYan struct {
	ID      bson.ObjectId `bson:"_id,omitempty"`
	Index   int           `bson:"index" json:"index"`
	Sex     string        `bson:"sex"`
	Fortune string        `bson:"fortune" json:"fortune"`
	TianJiu string        `bson:"tian_jiu" json:"tian_jiu"`
	Comment string        `bson:"comment" json:"comment"`
}

type Iterable interface {
	HasNext() bool
	Next() interface{}
	Reset()
	Add(v interface{})
	Size() int
	IteratorFunc(f IteratorFunc) error
}

//Character 字符
type Character struct {
	//URL            string //汉字地址
	IsCommonly     bool   `bson:"is_commonly"`     //是否为常用字
	Character      string `bson:"character"`       //字符
	Pinyin         string `bson:"pinyin"`          //拼音
	Radical        string `bson:"radical"`         //部首
	RadicalStrokes string `bson:"radical_strokes"` //部首笔画
	TotalStrokes   string `bson:"total_strokes"`   //总笔画
	KangxiStrokes  string `bson:"kangxi_strokes"`  //康熙笔画数
	Phonetic       string `bson:"phonetic"`        //注音
	Folk           `json:"folk"`
	Structure      `json:"structure"`
	Explain        `json:"explain"`
	Rhyme          `json:"rhyme"`
	Index          `json:"index"`
}

//Folk民俗参考
type Folk struct {
	CommonlyCharacters   string `bson:"commonly_characters"`    //是否为常用字
	NameScience          string `bson:"name_science"`           //姓名学
	FiveElementCharacter string `bson:"five_element_character"` //汉字五行
	GodBadMoral          string `bson:"god_bad_moral"`          //吉凶寓意
}

//Structure 字形结构
type Structure struct {
	DecompositionSearch string `bson:"decomposition_search"` //首尾分解查字
	StrokeNumber        string `bson:"stroke_number"`        //笔顺编号
	StrokeReadWrite     string `bson:"stroke_read_write"`    //笔顺读写
}

//Explain 康熙字典解释
type Explain struct {
	Intro  string `bson:"intro"`  //简介
	Detail string `bson:"detail"` //详情
}

//Rhyme 音韵参考
type Rhyme struct {
	GuangYun  string `bson:"guang_yun"` //广　韵
	Mandarin  string `bson:"mandarin"`  //国　语
	Cantonese string `bson:"cantonese"` //粤　语
}

//Index 索引参考
type Index struct {
	AncientWrite      string `bson:"ancient_write"`      //古文字诂林
	HometownTrain     string `bson:"hometown_train"`     //故训彙纂
	Explain           string `bson:"explain"`            //说文解字
	KangxiDictionary  string `bson:"kangxi_dictionary"`  //康熙字典
	ChineseDictionary string `bson:"chinese_dictionary"` //汉语字典
	Cihai             string `bson:"cihai"`              //辞　海  
}

type RadicalCharacter struct {
	*RootCharacter `json:"root_character"`
	Strokes        string `bson:"strokes"`
	Pinyin         string `bson:"pinyin"`
	Character      string `bson:"character"`
	URL            string `bson:"url"`
}

type RootCharacter struct {
	Strokes string `bson:"strokes"`
	Name    string `bson:"name"`
	URL     string `bson:"url"`
}
