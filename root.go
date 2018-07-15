package excavator

type RootFunc func(radical *RootCharacter) error

//RootRadical result root list
type Root struct {
	beforeIterator RootFunc
	iterator
	afterIterator RootFunc
	//wg     sync.WaitGroup
	URL    string
	Suffix string
}

type RootCharacter struct {
	Strokes string
	Name    string
	URL     string
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

func (root *Root) Iterator(f RootFunc) {
	root.Reset()
	for root.HasNext() {
		rc := root.Next().(*RootCharacter)
		if root.beforeIterator != nil {
			if err := root.beforeIterator(rc); err != nil {
				panic(err)
			}
		}
		if err := f(rc); err != nil {
			panic(err)
		}
		if root.afterIterator != nil {
			if err := root.afterIterator(rc); err != nil {
				panic(err)
			}
		}

	}
}
