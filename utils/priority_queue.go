package utils

import (
	"sort"
)

type PriorityQueue struct {
	array      []interface{}
	comparator func(i, j int) bool
}

func NewPriorityQueue(comparator func(pq *PriorityQueue) func(i, j int) bool) *PriorityQueue {
	pq := &PriorityQueue{
		array: make([]interface{}, 0),
	}
	pq.comparator = comparator(pq)
	return pq
}

func (pq *PriorityQueue) Get(i int) interface{} {
	return pq.array[i]
}

func (pq *PriorityQueue) Add(x interface{}, onAdd func(interface{})) {
	pq.array = append(pq.array, x)

	first := pq.array[0]
	sort.Slice(pq.array, pq.comparator)

	if first != pq.array[0] {
		onAdd(first)
	}
}

func (pq *PriorityQueue) Peek() interface{} {
	if len(pq.array) == 0 {
		return nil
	}
	elem := pq.array[0]
	return elem
}

func (pq *PriorityQueue) Pop() interface{} {
	if len(pq.array) == 0 {
		return nil
	}
	first, rest := pq.array[0], pq.array[1:]
	pq.array = rest
	return first
}

func (pq *PriorityQueue) Len() int {
	return len(pq.array)
}
