package excavator

type IteratorFunc func(v interface{}) error

type iterator struct {
	Max   int
	data  []interface{}
	index int
}

func NewIterator() *iterator {
	return &iterator{
		Max:   3,
		data:  nil,
		index: 0,
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

//Size iterator data size
func (i *iterator) Size() int {
	return len(i.data)
}

//Iterator an default iterator
func (i *iterator) IteratorFunc(f IteratorFunc) error {
	i.Reset()
	for i.HasNext() {
		if err := f(i.Next()); err != nil {
			return err
		}
	}
	return nil
}

func (i *iterator) Iterator() *iterator {
	return i
}

func (i *iterator) Data() []interface{} {
	return i.data
}

func process(iter *iterator, f IteratorFunc, cb chan bool) bool {
	if iter.HasNext() {
		go func(cb chan<- bool) {
			defer func() {
				cb <- true
			}()
			if err := f(iter.Next()); err != nil {
				return
			}
		}(cb)
		return true
	}
	return false
}

func (iter *iterator) ThreadIterator(f IteratorFunc) error {
	cb := make(chan bool, iter.Max)
	defer close(cb)
	iter.Reset()

	for i1 := 0; i1 < iter.Max; i1++ {
		if !process(iter, f, cb) {
			return nil
		}
	}
	for {
		select {
		case <-cb:
			if !process(iter, f, cb) {
				return nil
			}
		}
	}
	return nil
}
