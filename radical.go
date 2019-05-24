package excavator

type RadicalFunc func(rc *RadicalCharacter) error

//Radical
type Radical struct {
	root           *Root
	beforeIterator RadicalFunc
	iterator
}

type RadicalCharacter struct {
	*RootCharacter `json:"rootcharacter"`
	Strokes   string `json:"strokes"`
	Pinyin    string `json:"pinyin"`
	Character string `json:"character"`
	URL       string `json:"url"`
}

func (r *Radical) Add(rc *RadicalCharacter) {
	r.iterator.Add(rc)
}

func (r *Radical) SetRoot(root *Root) {
	r.root = root
}

func (r *Radical) IteratorFunc(f RadicalFunc) []*Character {
	var rad []*Character
	r.Reset()
	for r.HasNext() {
		rc := r.Next().(*RadicalCharacter)
		if r.beforeIterator != nil {
			if err := r.beforeIterator(rc); err != nil {
				log.Panic(err)
			}
		}
		rad = append(rad, getCharacterList(r.root, rc))
		if err := f(rc); err != nil {
			log.Panic(err)
			continue
		}
	}
	return rad
}

func (r *Radical) Character(character *RadicalCharacter) *Character {
	return getCharacterList(r.root, character)
}

func (r *Radical) SetBefore(f RadicalFunc) {
	r.beforeIterator = f
}
