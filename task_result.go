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
	"sync"
	"sync/atomic"
	"time"
)

// WaitAllTaskResults waits until all the tasks terminate.
func WaitAllTaskResults(results ...*TaskResult) {
	if len(results) == 0 {
		return
	}

	var wg sync.WaitGroup
	for _, result := range results {
		wg.Add(1)
		callback := result.callback
		result.callback = func(tr *TaskResult) {
			wg.Done()
			if callback != nil {
				callback(tr)
			}
		}
	}
	wg.Wait()
}

// TaskResult represents the result of the task.
type TaskResult struct {
	wid  uint64
	pool *TaskPool
	task Task
	args []interface{}
	schd chan struct{}
	done chan struct{}

	start time.Time
	cost  time.Duration
	res   interface{}
	err   error

	lock     sync.RWMutex
	callback func(*TaskResult)
}

// Release releases the task result, which will put it into the pool.
//
// It is optional, not mandatory. So, once calling it, you must not cache
// and use it again.
func (r *TaskResult) Release() { r.pool.putTaskResult(r) }
func (r *TaskResult) release() { *r = TaskResult{} }

func (r *TaskResult) setwid(wid uint64) {
	atomic.StoreUint64(&r.wid, wid)
	r.start = time.Now()
	close(r.schd)
}
func (r *TaskResult) reset(pool *TaskPool, task Task, args []interface{}) {
	r.task = task
	r.args = args
	r.pool = pool
	r.cost = 0
	r.res = nil
	r.err = nil
	r.schd = make(chan struct{})
	r.done = make(chan struct{})
}

func (r *TaskResult) setResult(result interface{}, err error) {
	r.err = err
	r.res = result
	r.cost = time.Since(r.start)
	close(r.done)

	r.lock.RLock()
	cb := r.callback
	r.lock.RUnlock()
	if cb != nil {
		cb(r)
	}
}

// Task returns the associatived task.
func (r *TaskResult) Task() Task { return r.task }

// Args returns all the arguments passed to the task.
func (r *TaskResult) Args() []interface{} { return r.args }

// Pool returns the task pool running the task that produces the result.
func (r *TaskResult) Pool() *TaskPool { return r.pool }

// WorkerID returns the id of the worker to run the task.
func (r *TaskResult) WorkerID() uint64 { return atomic.LoadUint64(&r.wid) }

// Done return a channel to indicate when the task has done.
func (r *TaskResult) Done() <-chan struct{} { return r.done }

// Wait waits until the task has done, which is equal to <-Done().
func (r *TaskResult) Wait() { <-r.done }

// WaitSchedule waits until the task is scheduled.
func (r *TaskResult) WaitSchedule() { <-r.schd }

// IsScheduled reports whether the task has been scheduled, that's,
// there is a worker to run the task, which is equal to "r.WorkerID() > 0".
func (r *TaskResult) IsScheduled() bool { return r.WorkerID() > 0 }

// IsDone reports whether the task has already done.
func (r *TaskResult) IsDone() bool {
	select {
	case <-r.done:
		return true
	default:
		return false
	}
}

// IsSuccess reports whether the task has done successfully,
// that's, the task didn't return the error.
//
// Notice: it will wait until the task has done.
func (r *TaskResult) IsSuccess() bool { return r.Error() == nil }

// Result returns the result returned by the task.
//
// Notice: it will wait until the task has done.
func (r *TaskResult) Result() interface{} { r.Wait(); return r.res }

// Error returns the error returned by the task.
//
// Notice: it will wait until the task has done.
func (r *TaskResult) Error() error { r.Wait(); return r.err }

// StartTime returns the time when the task is started.
func (r *TaskResult) StartTime() time.Time { return r.start }

// Duration returns the cost duration to run the task.
func (r *TaskResult) Duration() time.Duration { r.Wait(); return r.cost }

// SetCallback sets the callback function, and calls it when the task has done.
func (r *TaskResult) SetCallback(callback func(*TaskResult)) {
	if callback == nil {
		panic("TaskResult: the callback function must not be nil")
	}

	r.lock.Lock()
	r.callback = callback
	r.lock.Unlock()
	if r.IsDone() {
		callback(r)
	}
}
