package queue

import "testing"

func TestLen(t *testing.T) {
	list := ItemList{}
	if list.Len() != 0 {
		t.Error("The list should be empty")
	}
}

func TestPush(t *testing.T) {
	list := ItemList{}
	list.Push(NewItem("github.com/foo/bar", 20))
	if list.Len() != 1 {
		t.Error("There should be one item")
	}
}

func TestLess(t *testing.T) {
	list := ItemList{}
	item1 := NewItem("github.com/foo/bar", 20)
	item2 := NewItem("github.com/bar/foo", 50)
	list.Push(item1)
	list.Push(item2)
	if !list.Less(1, 0) {
		t.Error("The second item should have a higher priority")
	}
}

func TestSwap(t *testing.T) {
	list := ItemList{}
	item1 := NewItem("github.com/foo/bar", 20)
	item2 := NewItem("github.com/bar/foo", 50)
	list.Push(item1)
	list.Push(item2)
	list.Swap(1, 0)
	if list.Pop().(*Item).value != "github.com/foo/bar" {
		t.Error("Unexpected last item")
	}
}
