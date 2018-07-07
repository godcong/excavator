package excavator

//RootRadical result root list
type Root struct {
	Radicals map[int]Radical
}

type Radical struct {
	Strokes string
	Name    string
	URL     string
}

func NewRoot() *Root {
	return &Root{
		Radicals: make(map[int]Radical),
	}
}

//func (r *Root)Add(key,v)  {
//
//}