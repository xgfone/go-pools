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
)

// PanicError converts an panic to error.
type PanicError struct{ Err interface{} }

func (pe PanicError) Error() string { return fmt.Sprintf("panic: %v", pe.Err) }

type taskCtx struct {
	ctx  context.Context
	task Task
	args []interface{}

	result *TaskResult
}

type worker struct {
	id    uint64
	pool  *TaskPool
	queue chan taskCtx
	exit  chan struct{}
	done  chan struct{}
}

func newWorker(id uint64, pool *TaskPool) *worker {
	return &worker{
		id:    id,
		pool:  pool,
		exit:  make(chan struct{}),
		done:  make(chan struct{}),
		queue: make(chan taskCtx),
	}
}

func (w *worker) handlePanic(r *TaskResult) {
	if err := recover(); err != nil && r != nil {
		if handler := w.pool.getPanicHandler(); handler != nil {
			r.setResult(nil, PanicError{Err: err})
			handler(r)
		}
	}
}

func (w *worker) handleTask(tc taskCtx) {
	defer w.handlePanic(tc.result)
	result, err := tc.task.Run(tc.ctx, tc.args...)
	if tc.result != nil {
		tc.result.setResult(result, err)
	}
	w.pool.taskIsDone(tc.result)
}

func (w *worker) Submit(tc taskCtx) bool {
	if tc.result != nil {
		tc.result.setwid(w.id)
	}

	select {
	case <-w.exit:
		return false
	case w.queue <- tc:
		return true
	}
}

func (w *worker) Wait() { <-w.done }
func (w *worker) Stop() {
	select {
	case <-w.exit:
	default:
		close(w.exit)
	}
}

func (w *worker) Start() {
	defer recover()
	w.pool.activateWorker(w)

	for {
		select {
		case <-w.exit:
			// Execute the last task to prevent the task being gone.
			select {
			case tc := <-w.queue:
				w.handleTask(tc)
			default:
			}

			close(w.done)
			return

		case tc := <-w.queue:
			w.handleTask(tc)
			select {
			case <-w.exit:
				close(w.done)
				return

			default:
				w.pool.activateWorker(w)
			}
		}
	}
}
