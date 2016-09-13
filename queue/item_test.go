package queue

import "testing"

func TestNewItem(t *testing.T) {
	value := "github.com/foo/bar"
	priority := 20
	item := NewItem(value, priority)
	if item.value != value {
		t.Error("Unexpected value")
	}
	if item.priority != priority {
		t.Error("Unexpected priority")
	}
	if item.hash == 0 {
		t.Error("The hash should not be empty")
	}
}
