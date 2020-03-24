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

func (pq *PriorityQueue) Add(x interface{}) {
	pq.array = append(pq.array, x)
	sort.Slice(pq.array, pq.comparator)
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

func (pq *PriorityQueue) Array() []interface{} {
	return pq.array
}
