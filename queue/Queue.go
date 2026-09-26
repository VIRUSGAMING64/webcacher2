package queue

import (
	"sync"
	"sync/atomic"
)

type QueueObj struct {
	Count int    `json:"count"`
	Url   string `json:"url"`
}

type Queue struct {
	mtx     sync.Mutex
	workers int
	running atomic.Int32
}

func NewQueue() *Queue {
	q := Queue{
		workers: 8,
		mtx:     sync.Mutex{},
		running: atomic.Int32{},
	}
	return &q
}

var tQueue *Queue = NewQueue()
