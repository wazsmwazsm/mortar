package mortar

import (
	"errors"
	"sync/atomic"
)

var (
	// ErrInvalidPoolCap return if pool size <= 0
	ErrInvalidPoolCap = errors.New("invalid pool cap")
)

// Task task to-do
type Task func()

// Pool task pool
type Pool struct {
	capacity       uint64
	runningWorkers uint64
	taskC          chan Task
}

// NewPool init pool
func NewPool(capacity uint64) (*Pool, error) {
	if capacity <= 0 {
		return nil, ErrInvalidPoolCap
	}
	return &Pool{
		capacity: capacity,
		taskC:    make(chan Task, capacity),
	}, nil
}

// GetCap get capacity
func (p *Pool) GetCap() uint64 {
	return atomic.LoadUint64(&p.capacity)
}

// GetRunningWorkers get running workers
func (p *Pool) GetRunningWorkers() uint64 {
	return atomic.LoadUint64(&p.runningWorkers)
}
