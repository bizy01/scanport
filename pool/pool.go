package pool

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"
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


type Task struct {
	Handler func(v ...interface{})
	Params []interface{}
}

type Pool struct {
	capacity uint64
	runningWorkNums uint64
	state int
	taskC chan *Task
	sync.Mutex
}

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


// 获取容量
func (p *Pool) GetCap() uint64 {
	return p.capacity
}

func (p *Pool) GetRunningWorks() uint64 {
	return atomic.LoadUint64(&p.runningWorkNums)
}

// put task
func (p *Pool) Put(task *Task) error {
	if p.getState() == STOPED {
		return ErrPoolAlreadyClosed
	}

	// safe run worker
	p.Lock()
	if p.GetRunningWorks() < p.GetCap() {
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

func (p *Pool) incr() {
	atomic.AddUint64(&p.runningWorkNums, 1)
}

func (p *Pool) decr() {
	atomic.AddUint64(&p.runningWorkNums, ^uint64(0))
}

// run
func (p *Pool) run() {
	p.incr()

	go func() {
		defer p.decr()
		for {
			select {
			case task, ok := <- p.taskC:
				if !ok {
					return
				}

				task.Handler(task.Params...)
			}
		}
	}()
}

func (p *Pool) getState() int {
	p.Lock()
	defer p.Unlock()

	return p.state
}

func (p *Pool) setState(state int) {
	p.Lock()
	defer p.Unlock()

	p.state = state
}

func (p *Pool) close() {
	p.Lock()
	defer p.Unlock()

	close(p.taskC)
}

func (p *Pool) Close() {
	if p.getState() == STOPED {
		return
	}

	p.setState(STOPED)

	for len(p.taskC) > 0 {
		time.Sleep(1e6)
	}

	p.close()
}