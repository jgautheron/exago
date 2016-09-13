package queue

// ItemList is a representation of the queue.
// It contains all items that remain to be processed.
type ItemList []*Item

func (pq ItemList) Len() int { return len(pq) }

func (pq ItemList) Less(i, j int) bool {
	return pq[i].priority > pq[j].priority
}

func (pq ItemList) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *ItemList) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *ItemList) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.index = -1
	*pq = old[0 : n-1]
	return item
}
