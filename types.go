package excavator

//RootRadical result root list
type Root struct {
	iterator int
	Radicals []*Radical
	URL      string
}

//Radical
type Radical struct {
	Strokes string
	Name    string
	URL     string
}

type Character struct {
	Character      string //字符
	Pinyin         string //拼音
	Radical        string //部首
	RadicalStrokes int    //部首笔画
	KangxiStrokes  int    //康熙笔画数
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
		iterator: 0,
		Radicals: make([]*Radical,0),
		URL:      url,
	}
}

//Add add radical
func (r *Root) Add(rd *Radical) {
	r.Radicals = append(r.Radicals, rd)
}

//HasNext check next
func (r *Root) HasNext() bool {
	return r.iterator < len(r.Radicals)
}

//Next get next
func (r *Root) Next() *Radical {
	defer func() {
		r.iterator++
	}()
	if r.iterator < len(r.Radicals) {
		return r.Radicals[r.iterator]
	}

	return nil
}

//Reset reset index
func (r *Root) Reset() {
	r.iterator = 0
}

//GetList get an list from web
func (root *Root) GetList(s string) {
	getRootList(root, s)
	getRedicalList(root)
}
