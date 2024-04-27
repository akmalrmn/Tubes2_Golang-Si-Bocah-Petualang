package priorityqueue

import (
	"be/pkg/tree"
	"container/heap"
	"sync"
)

type Item struct {
	Value    *tree.Node
	priority int
	index    int
}

type PriorityQueue []*Item

func (pq *PriorityQueue) Len() int { return len(*pq) }

func (pq *PriorityQueue) Less(i, j int) bool {
	return (*pq)[i].priority < (*pq)[j].priority
}

func (pq *PriorityQueue) Swap(i, j int) {
	(*pq)[i], (*pq)[j] = (*pq)[j], (*pq)[i]
	(*pq)[i].index = i
	(*pq)[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*pq = old[0 : n-1]
	return item
}

type PriorityChannel struct {
	pq     PriorityQueue
	c      chan *Item
	mu     sync.Mutex
	closed bool
}

func NewPriorityChannel() *PriorityChannel {
	pc := &PriorityChannel{
		c: make(chan *Item),
	}
	heap.Init(&pc.pq)
	go pc.run()
	return pc
}

func (pc *PriorityChannel) run() {
	for {
		pc.mu.Lock()
		for pc.pq.Len() == 0 && !pc.closed {
			pc.mu.Unlock()
			pc.mu.Lock()
		}
		if pc.pq.Len() == 0 && pc.closed {
			pc.mu.Unlock()
			close(pc.c)
			return
		}
		item := heap.Pop(&pc.pq).(*Item)
		pc.mu.Unlock()
		pc.c <- item
	}
}

func (pc *PriorityChannel) Add(item *tree.Node) {
	pc.mu.Lock()
	defer pc.mu.Unlock()
	if pc.closed {
		return
	}
	treeItem := &Item{
		Value:    item,
		priority: item.Depth,
		index: item.Depth,
	}
	heap.Push(&pc.pq, treeItem)
}

func (pc *PriorityChannel) Close() {
	pc.mu.Lock()
	defer pc.mu.Unlock()
	pc.closed = true
}

func (pc *PriorityChannel) C() <-chan *Item {
	return pc.c
}