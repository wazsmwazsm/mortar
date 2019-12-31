package mortar

import (
	"context"
	"errors"
	"log"
	"sync/atomic"
)

var (
	// ErrInvalidPoolCap return if pool size <= 0
	ErrInvalidPoolCap = errors.New("invalid pool cap")
	// ErrPoolAlreadyClosed put task but pool already closed
	ErrPoolAlreadyClosed = errors.New("pool already closed")
)

const (
	// RUNNING pool is running
	RUNNING = 1
	// STOPED pool is stoped
	STOPED = 0
)

// Task task to-do
type Task struct {
	Handler func(v ...interface{})
	Params  []interface{}
}

// Pool task pool
type Pool struct {
	capacity       uint64
	runningWorkers uint64
	state          int64
	taskC          chan *Task
	closeC         chan bool
	PanicHandler   func(interface{})
}

// NewPool init pool
func NewPool(capacity uint64) (*Pool, error) {
	if capacity <= 0 {
		return nil, ErrInvalidPoolCap
	}
	return &Pool{
		capacity: capacity,
		state:    RUNNING,
		taskC:    make(chan *Task, capacity),
		closeC:   make(chan bool),
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

func (p *Pool) incRunning() {
	atomic.AddUint64(&p.runningWorkers, 1)
}

func (p *Pool) decRunning() {
	atomic.AddUint64(&p.runningWorkers, ^uint64(0))
}

// Put put a task to pool
func (p *Pool) Put(ctx context.Context, task *Task) error {

	if p.state == STOPED {
		return ErrPoolAlreadyClosed
	}

	if p.GetRunningWorkers() < p.GetCap() {
		p.run(ctx)
	}

	p.taskC <- task

	return nil
}

func (p *Pool) run(ctx context.Context) {
	p.incRunning()

	go func(ctx context.Context) {
		defer func() {
			p.decRunning()
			if r := recover(); r != nil {
				if p.PanicHandler != nil {
					p.PanicHandler(r)
				} else {
					log.Printf("Worker panic: %s\n", r)
				}
			}
		}()

		for {
			select {
			case task, ok := <-p.taskC:
				if !ok {
					return
				}
				task.Handler(task.Params...)
			case <-ctx.Done():
				return
			case <-p.closeC:
				return
			}
		}
	}(ctx)
}

// Close close pool graceful
func (p *Pool) Close() {
	p.state = STOPED // stop put task

	for len(p.taskC) > 0 { // wait all task be consumed
	}

	p.closeC <- true
	close(p.taskC)
}
