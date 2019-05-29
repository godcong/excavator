package excavator

// IteratorFunc ...
type IteratorFunc func(v interface{}) error

type iterator struct {
	Max   int
	data  []interface{}
	index int
}

// Iterator ...
type Iterator interface {
	HasNext() bool
	Next() interface{}
	Reset()
	Add(v interface{})
	Size() int
	IteratorFunc(f IteratorFunc) error
	Iterator() Iterator
	ThreadIterator(f IteratorFunc) error
	Data() []interface{}
}

// NewIterator ...
func NewIterator() Iterator {
	return &iterator{
		Max:   3,
		data:  nil,
		index: 0,
	}
}

//HasNext check next
func (iter *iterator) HasNext() bool {
	return iter.index < len(iter.data)
}

//Next get next
func (iter *iterator) Next() interface{} {
	defer func() {
		iter.index++
	}()
	if iter.index < len(iter.data) {
		return iter.data[iter.index]
	}

	return nil
}

//Reset reset index
func (iter *iterator) Reset() {
	iter.index = 0
}

//Add add radical
func (iter *iterator) Add(v interface{}) {
	iter.data = append(iter.data, v)
}

//Size iterator data size
func (iter *iterator) Size() int {
	return len(iter.data)
}

//Iterator an default iterator
func (iter *iterator) IteratorFunc(f IteratorFunc) error {
	iter.Reset()
	for iter.HasNext() {
		if err := f(iter.Next()); err != nil {
			return err
		}
	}
	return nil
}

// Iterator ...
func (iter *iterator) Iterator() Iterator {
	return iter
}

// Data ...
func (iter *iterator) Data() []interface{} {
	return iter.data
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

// ThreadIterator ...
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
}
