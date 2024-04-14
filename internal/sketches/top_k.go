package sketches

import (
	"container/heap"
	"math"
	"sort"
)

// Top k Item
type TopKItem struct {
	Key       string
	Frequency uint32
}

type TopKItems []TopKItem

func (e TopKItems) Len() int           { return len(e) }
func (e TopKItems) Less(i, j int) bool { return e[i].Frequency < e[j].Frequency }
func (e TopKItems) Swap(i, j int)      { e[i], e[j] = e[j], e[i] }

// TopKSketch struct
type TopKSketch struct {
	cms *CountMinSketch
	pq  TopKPriorityQueue
}

// TopK struct using count min sketch and priority queue
type TopKPriorityQueue []*TopKItem

// ensure that TopKPriorityQueue implements the heap interface
var _ heap.Interface = (*TopKPriorityQueue)(nil)

// Len method for TopKPriorityQueue
func (pq TopKPriorityQueue) Len() int { return len(pq) }
func (pq TopKPriorityQueue) Less(i, j int) bool {
	return pq[i].Frequency < pq[j].Frequency
}
func (pq TopKPriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}
func (pq *TopKPriorityQueue) Push(x any) {
	item := x.(*TopKItem)
	*pq = append(*pq, item)
}
func (pq *TopKPriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	*pq = old[0 : n-1]
	return item
}

// Find returns the index of a given item in the priority queue
// so that it can be updated in the heap. returns -1 if the item is not found
func (pq TopKPriorityQueue) Find(item string) int {
	for i, v := range pq {
		if v.Key == item {
			return i
		}
	}
	return -1
}

// return the min frequency in the priority queue
func (pq TopKPriorityQueue) Min() uint32 {
	return pq[0].Frequency
}
func (pq *TopKPriorityQueue) update(index int, frequency uint32) {
	(*pq)[index].Frequency = frequency
	heap.Fix(pq, index)
}

func NewTopK(k int) *TopKSketch {
	// Using some findings from Heavy Keepers algorithm for releation between width/depth and k and setting some min values
	width := float64(k) * math.Log(float64(k))
	depth := math.Log(float64(k))

	if width < 100 {
		width = 100
	}
	if depth < 3 {
		depth = 3
	}

	pq := make(TopKPriorityQueue, 0)

	for _ = range k {
		pq = append(pq, &TopKItem{})
	}
	return &TopKSketch{
		cms: NewCountMinSketch(uint32(width), uint32(depth)),
		pq:  pq,
	}
}

func (t *TopKSketch) Reset() {
	t.cms.Reset()
	for i := range t.pq {
		t.pq[i] = &TopKItem{}
	}
}

// Add element in cms, also update the priority queue
// Redis style top add, if an element was dropped from the TopK list, Nil reply otherwise..
func (t *TopKSketch) Add(key string) (expelledItem string) {
	t.cms.Add(key)
	freq := t.cms.Count(key)

	minHeap := t.pq.Min()
	if freq >= minHeap {
		index := t.pq.Find(key)
		if index > -1 {
			t.pq.update(index, freq)
			heap.Fix(&t.pq, index)
		} else {
			expelledItem = t.pq[0].Key
			t.pq[0].Frequency = freq
			t.pq[0].Key = key
			heap.Fix(&t.pq, 0)
			return expelledItem
		}
	}
	return ""
}

func (t *TopKSketch) List() []TopKItem {
	pqCopy := make([]TopKItem, len(t.pq))
	for i, item := range t.pq {
		pqCopy[i] = (*item)
	}
	sort.Stable(sort.Reverse(TopKItems(pqCopy)))

	// remove empty values
	r := len(t.pq)
	for ; r > 0; r-- {
		if pqCopy[r-1].Frequency > 0 {
			break
		}
	}
	return pqCopy[:r]
}
func (t *TopKSketch) Query(key string) uint32 {

	for i := range t.pq {
		if t.pq[i].Key == key {
			return t.pq[i].Frequency
		}
	}
	return 0
}
