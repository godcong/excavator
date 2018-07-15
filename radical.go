package excavator

type RadicalFunc func(radical *RadicalCharacter) error

//Radical
type Radical struct {
	root *Root
	iterator
}

type RadicalCharacter struct {
	*RootCharacter
	Strokes   string
	Pinyin    string
	Character string
	URL       string
}

func (r *Radical) Add(rc *RadicalCharacter) {
	r.iterator.Add(rc)
}

func (r *Radical) SetRoot(root *Root) {
	r.root = root
}
