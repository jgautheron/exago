package queue

import (
	"fmt"
	"hash/crc32"
	"time"
)

// Item is the definition of a message to be processed.
// The index is maintained by the heap, do not set it manually!
type Item struct {
	hash     uint32
	value    string
	priority int
	index    int
}

// NewItem creates a new message with the given value and priority.
// The hash is used as UUID.
func NewItem(value string, priority int) *Item {
	// CRC32 will do for now...
	hash := crc32.ChecksumIEEE([]byte(fmt.Sprintf(
		"%s-%s",
		value,
		time.Now(),
	)))

	return &Item{
		value:    value,
		priority: priority,
		hash:     hash,
	}
}
