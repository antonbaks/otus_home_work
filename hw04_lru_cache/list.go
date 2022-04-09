package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	len   int
	front *ListItem
	back  *ListItem
}

func NewList() List {
	return new(list)
}

func (l list) Len() int {
	return l.len
}

func (l list) Front() *ListItem {
	return l.front
}

func (l list) Back() *ListItem {
	return l.back
}

func (l *list) Remove(i *ListItem) {
	l.len--

	l.removeListItem(i)
}

func (l *list) removeListItem(i *ListItem) {
	if i.Prev == nil {
		l.removeFront(i)

		return
	}

	if i.Next == nil {
		l.removeBack(i)

		return
	}

	i.Prev.Next = i.Next
}

func (l *list) removeFront(oldFrontListItem *ListItem) {
	newFrontListItem := oldFrontListItem.Next
	newFrontListItem.Prev = nil
	l.front = newFrontListItem
}

func (l *list) removeBack(oldBackListItem *ListItem) {
	newFrontListItem := oldBackListItem.Prev
	newFrontListItem.Next = nil
	l.back = newFrontListItem
}

func (l *list) MoveToFront(i *ListItem) {
	l.removeListItem(i)
	l.addNewFrontListItem(i)
}

func (l *list) PushFront(v interface{}) *ListItem {
	newFrontListItem := ListItem{
		v,
		nil,
		nil,
	}

	l.addNewFrontListItem(&newFrontListItem)

	l.len++

	return l.Front()
}

func (l *list) addNewFrontListItem(i *ListItem) {
	front := l.Front()
	if front != nil {
		front.Prev = i
		i.Next = front
	} else {
		l.bindWithBack(i)
	}

	l.front = i
}

func (l *list) bindWithBack(newFrontListItem *ListItem) {
	back := l.Back()

	if back == nil {
		l.back = newFrontListItem

		return
	}

	for back.Prev != nil {
		back = back.Prev
	}

	back.Prev = newFrontListItem
	newFrontListItem.Next = back
}

func (l *list) PushBack(v interface{}) *ListItem {
	newBackListItem := ListItem{
		v,
		nil,
		nil,
	}

	back := l.Back()
	if back != nil {
		back.Next = &newBackListItem
		newBackListItem.Prev = back
	} else {
		l.bindWithFront(&newBackListItem)
	}

	l.back = &newBackListItem

	l.len++

	return l.Back()
}

func (l *list) bindWithFront(backListItem *ListItem) {
	front := l.Front()

	if front == nil {
		l.front = backListItem

		return
	}

	for front.Next != nil {
		front = front.Next
	}

	front.Next = backListItem
	backListItem.Prev = front
}
