package excavator

// RadicalFunc ...
type RadicalFunc func(rc *RadicalCharacter) error

// Radical ...
type Radical struct {
	root           *Root
	beforeIterator RadicalFunc
	iterator
}

// RadicalCharacter ...
type RadicalCharacter struct {
	*RootCharacter `json:"rootcharacter"`
	Strokes        string `json:"strokes"`
	Pinyin         string `json:"pinyin"`
	Character      string `json:"character"`
	URL            string `json:"url"`
}

// Add ...
func (r *Radical) Add(rc *RadicalCharacter) {
	r.iterator.Add(rc)
}

// SetRoot ...
func (r *Radical) SetRoot(root *Root) {
	r.root = root
}

// IteratorFunc ...
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

// Character ...
func (r *Radical) Character(character *RadicalCharacter) *Character {
	return getCharacterList(r.root, character)
}

// SetBefore ...
func (r *Radical) SetBefore(f RadicalFunc) {
	r.beforeIterator = f
}
