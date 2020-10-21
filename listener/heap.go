package listener

import (
	"container/heap"
	"sync"
	"time"

	"beanbot/state"
)

const (
	MESSAGE_QUEUE_FLUSH_INTERVAL = 50
)

type MessageHeap struct {
	s sync.RWMutex
	m []*state.MessageData
}

func (m *MessageHeap) Len() int {
	return len(m.m)
}

func (m *MessageHeap) Less(i, j int) bool {
	return m.m[i].Timestamp.Before(m.m[j].Timestamp)
}

func (m *MessageHeap) Swap(i, j int) {
	m.m[i], m.m[j] = m.m[j], m.m[i]
}

func (m *MessageHeap) Push(x interface{}) {
	item := x.(*state.MessageData)
	m.m = append(m.m, item)
}

func (m *MessageHeap) Pop() interface{} {
	old := m.m
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	m.m = old[0 : n-1]
	return item
}

func (m *MessageHeap) ShouldPop() bool {
	return m.Len() != 0 && m.m[0].Timestamp.Before(time.Now().Add(time.Millisecond*-MESSAGE_QUEUE_FLUSH_INTERVAL))
}

func (m *MessageHeap) Flush() (flushed []*state.MessageData) {
	for m.ShouldPop() {
		v := heap.Pop(m)
		flushed = append(flushed, v.(*state.MessageData))
	}
	return
}
