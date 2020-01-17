package mortar

import (
	"errors"
	"log"
	"sync"
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
	PanicHandler   func(interface{})
	sync.Mutex
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
	}, nil
}

// GetCap get capacity
func (p *Pool) GetCap() uint64 {
	return p.capacity
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
func (p *Pool) Put(task *Task) error {

	if p.getState() == STOPED {
		return ErrPoolAlreadyClosed
	}

	// safe run worker
	p.Lock()
	if p.GetRunningWorkers() < p.GetCap() {
		p.run()
	}
	p.Unlock()

	// send task safe
	p.Lock()
	if p.state == RUNNING {
		p.taskC <- task
	}
	p.Unlock()

	return nil
}

func (p *Pool) run() {
	p.incRunning()

	go func() {
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
			}
		}
	}()
}

func (p *Pool) getState() int64 {
	p.Lock()
	defer p.Unlock()

	return p.state
}

func (p *Pool) setState(state int64) {
	p.Lock()
	defer p.Unlock()

	p.state = state
}

// close safe
func (p *Pool) close() {
	p.Lock()
	defer p.Unlock()

	close(p.taskC)
}

// Close close pool graceful
func (p *Pool) Close() {

	if p.getState() == STOPED {
		return
	}

	p.setState(STOPED) // stop put task

	for len(p.taskC) > 0 { // wait all task be consumed
	}

	p.close()
}
