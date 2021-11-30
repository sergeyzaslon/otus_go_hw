package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
	Clear()
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	first  *ListItem
	last   *ListItem
	length int
}

func (l *list) Len() int {
	return l.length
}

func (l *list) Front() *ListItem {
	return l.first
}

func (l *list) Back() *ListItem {
	return l.last
}

func (l *list) PushFront(v interface{}) *ListItem {
	l.length++
	li := &ListItem{
		Value: v,
		Next:  l.first,
		Prev:  nil,
	}
	if l.first != nil {
		l.first.Prev = li
	}
	l.first = li
	if l.last == nil {
		l.last = li
	}
	return li
}

func (l *list) PushBack(v interface{}) *ListItem {
	l.length++
	li := &ListItem{
		Value: v,
		Next:  nil,
		Prev:  l.last,
	}
	if l.last != nil {
		l.last.Next = li
	}
	l.last = li
	return li
}

func (l *list) Remove(i *ListItem) {
	l.length--
	switch {
	case i.Next == nil:
		{ // последний, потому как нет ссылки на следующий
			l.last = i.Prev
			i.Prev.Next = nil // становиться последним потому не на кого ссылаться
		}
	case i.Prev == nil:
		{ // удаляем первый
			l.first = i.Next
			i.Prev = nil
		}
	default:
		{ // не первй и не последний
			i.Prev.Next = i.Next
			i.Next.Prev = i.Prev
		}
	}
}

func (l *list) MoveToFront(i *ListItem) {
	if l.first != i {
		if i.Prev != nil {
			i.Prev.Next = i.Next
		}
		if i.Next != nil {
			i.Next.Prev = i.Prev
		}
		i.Prev = nil
		i.Next = l.first
		l.first.Prev = i
		l.first = i
	}
}

func (l *list) Clear() {
	l.first = nil
	l.last = nil
	l.length = 0
}

func NewList() List {
	return new(list)
}
