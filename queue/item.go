package queue

import (
	"fmt"
	"hash/crc32"
	"time"
)

type Item struct {
	hash     uint32
	value    string
	priority int
	index    int
}

func (i *Item) Hash() uint32 {
	i.hash = crc32.ChecksumIEEE([]byte(fmt.Sprintf(
		"%s-%s",
		i.value,
		time.Now(),
	)))
	return i.hash
}
