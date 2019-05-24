package excavator
type RootFunc func(rc *RootCharacter) error

//RootRadical result root list
type Root struct {
	beforeIterator RootFunc
	iterator
	URL    string
	Suffix string
}

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

func (root *Root) Self() *Root {
	return getRootList(root, root.Suffix)
}

func (root *Root) IteratorFunc(f RootFunc) []*Radical {
	var rad []*Radical
	root.Reset()
	for root.HasNext() {
		rc := root.Next().(*RootCharacter)
		if root.beforeIterator != nil {
			if err := root.beforeIterator(rc); err != nil {
				log.Panic(err)
				continue
			}
		}
		rad = append(rad, getRedicalList(root, rc))
		if err := f(rc); err != nil {
			log.Panic(err)
			continue
		}
	}
	return rad
}

func (root *Root) Radical(character *RootCharacter) *Radical {
	return getRedicalList(root, character)
}

func (root *Root) SetBefore(r RootFunc) {
	root.beforeIterator = r
}
