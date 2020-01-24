// Copyright 2019 xgfone
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package pools

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
)

// Task represents a task to be run.
type Task interface {
	Run(ctx context.Context, args ...interface{}) (interface{}, error)
}

// TaskFunc represents a function task.
type TaskFunc func(ctx context.Context, args ...interface{}) (interface{}, error)

// Run implements the Task inerface.
func (t TaskFunc) Run(ctx context.Context, args ...interface{}) (interface{}, error) {
	return t(ctx, args...)
}

// TaskStat represents the statistics of the task pool.
type TaskStat struct {
	Worker  uint64 `json:"worker"`  // The number of all the workers.
	Pending uint64 `json:"pending"` // The number of the tasks that will be run.
	Running uint64 `json:"running"` // The number of the tasks that is running.
	Done    uint64 `json:"done"`    // The total number of the tasks that have done.
}

func (ts TaskStat) String() string {
	return fmt.Sprintf("TaskPool(worker=%d, pending=%d, running=%d, done=%d)",
		ts.Worker, ts.Pending, ts.Running, ts.Done)
}

// TaskPool is an goroutine pool to run the task.
type TaskPool struct {
	done    chan struct{}
	exit    chan struct{}
	tasks   chan taskCtx
	handler func(*TaskResult)
	results sync.Pool
	hasdone uint64

	workers map[uint64]*worker
	wchan   chan *worker
}

// NewTaskPool returns a new task pool with the goroutine of the worker num.
//
// bufferSize is the size of the buffer of the task to be submitted,
// which is equal to workerNum by default.
func NewTaskPool(workerNum int, bufferSize ...int) *TaskPool {
	bsize := workerNum
	if len(bufferSize) > 0 {
		bsize = bufferSize[0]
	}

	if workerNum <= 0 {
		panic("TaskPool: the number of the task workers must be a positive integer")
	}
	if bsize <= 0 {
		panic("TaskPool: the size of the task buffer must be a positive integer")
	}

	pool := &TaskPool{
		done:    make(chan struct{}),
		exit:    make(chan struct{}),
		tasks:   make(chan taskCtx, bsize),
		wchan:   make(chan *worker, workerNum),
		workers: make(map[uint64]*worker, workerNum),
	}
	pool.addWorkers(workerNum)
	pool.results.New = func() interface{} { return &TaskResult{pool: pool} }
	go pool.run()
	return pool
}

func (p *TaskPool) getPanicHandler() func(*TaskResult) { return p.handler }
func (p *TaskPool) getTaskResult() *TaskResult         { return p.results.Get().(*TaskResult) }
func (p *TaskPool) putTaskResult(r *TaskResult)        { r.release(); p.results.Put(r) }
func (p *TaskPool) terminateWorker(w *worker)          { atomic.AddUint64(&p.hasdone, 1) }
func (p *TaskPool) activateWorker(w *worker)           { p.wchan <- w }

func (p *TaskPool) addWorkers(n int) {
	for i := 1; i <= n; i++ {
		wid := uint64(i)
		worker := newWorker(wid, p)
		p.workers[wid] = worker
		go worker.Start()
	}
}

func (p *TaskPool) run() {
	for {
		select {
		case <-p.exit:
			p.stopWorkers()
			return
		case r := <-p.tasks:
			var submit bool
			for !submit {
				select {
				case <-p.exit:
					p.stopWorkers()
					return
				case w := <-p.wchan:
					submit = w.Submit(r)
				}
			}
		}
	}
}

func (p *TaskPool) stopWorkers() {
	// Stop all the workers.
	for _, worker := range p.workers {
		worker.Stop()
	}

	// Wait until all the workers exit.
	for _, worker := range p.workers {
		worker.Wait()
	}

	// Notice someone that all the workers have exited.
	close(p.done)
}

func (p *TaskPool) stop() { close(p.exit) }

// Stop stops all the task in the pool and releases all the goroutines.
//
// Notice: it will wait until all the workers exit.
func (p *TaskPool) Stop() { p.Shutdown(context.Background()) }

// Shutdown is the same as Stop, but it waits until all the workers exit
// or the context is canceled.
//
// Notice: if context is nil, it will return immediately, and not wait that
// all the workers exit.
func (p *TaskPool) Shutdown(ctx context.Context) {
	p.stop()
	if ctx != nil {
		select {
		case <-ctx.Done():
		case <-p.done:
		}
	}
}

// Wait waits until the pool is stopped and all the workers exit.
func (p *TaskPool) Wait() { <-p.done }

// SetPanicHandler sets the panic handler.
//
// When panicking, it will wrap it and set the error of TaskResult to it.
//
// It does nothing by default.
func (p *TaskPool) SetPanicHandler(h func(*TaskResult)) { p.handler = h }

// GetPoolSize returns the size of the pool.
func (p *TaskPool) GetPoolSize() int { return len(p.workers) }

// TaskStat returns the statistics of the task pool.
func (p *TaskPool) TaskStat() (ts TaskStat) {
	ts.Worker = uint64(len(p.workers))
	ts.Pending = uint64(len(p.tasks))
	ts.Running = ts.Worker - uint64(len(p.wchan))
	ts.Done = atomic.LoadUint64(&p.hasdone)
	return
}

func (p *TaskPool) addTask(ctx context.Context, task Task, args ...interface{}) *TaskResult {
	r := p.getTaskResult()
	r.reset(p, task, args)
	p.tasks <- taskCtx{ctx: ctx, task: task, args: args, result: r}
	return r
}

func (p *TaskPool) addTask2(ctx context.Context, task Task, args ...interface{}) {
	p.tasks <- taskCtx{ctx: ctx, task: task, args: args}
}

// RunTask starts a task to run.
func (p *TaskPool) RunTask(task Task, args ...interface{}) {
	p.addTask2(context.Background(), task, args...)
}

// RunTaskFunc starts a task function to run.
func (p *TaskPool) RunTaskFunc(f TaskFunc, args ...interface{}) {
	p.addTask2(context.Background(), f, args...)
}

// RunTaskWithContext starts a task to run with the context.
func (p *TaskPool) RunTaskWithContext(ctx context.Context, task Task, args ...interface{}) {
	p.addTask2(ctx, task, args...)
}

// RunTaskFuncWithContext starts a task function to run with the context.
func (p *TaskPool) RunTaskFuncWithContext(ctx context.Context, f TaskFunc, args ...interface{}) {
	p.addTask2(ctx, f, args...)
}

// RunTaskWithResult starts a task to run and returns the result.
func (p *TaskPool) RunTaskWithResult(task Task, args ...interface{}) *TaskResult {
	return p.addTask(context.Background(), task, args...)
}

// RunTaskFuncWithResult starts a task function to run and returns the result.
func (p *TaskPool) RunTaskFuncWithResult(f TaskFunc, args ...interface{}) *TaskResult {
	return p.addTask(context.Background(), f, args...)
}

// RunTaskWithResultAndContext starts a task to run with the context
// and returns the result.
func (p *TaskPool) RunTaskWithResultAndContext(ctx context.Context, task Task, args ...interface{}) *TaskResult {
	return p.addTask(ctx, task, args...)
}

// RunTaskFuncWithResultAndContext starts a task function to run with the context
// and returns the result.
func (p *TaskPool) RunTaskFuncWithResultAndContext(ctx context.Context, f TaskFunc, args ...interface{}) *TaskResult {
	return p.addTask(ctx, f, args...)
}
