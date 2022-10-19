// Copyright 2022 xgfone
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

import "sync"

// Object is used to enclose T.
type Object[T any] struct {
	Object T
	pool   *Pool[T]
}

// Release releases the object into the original pool.
func (o *Object[T]) Release() {
	if o != nil && o.pool != nil {
		o.pool.Put(o)
	}
}

// Pool is the object pool to allocate an object.
type Pool[T any] struct {
	pool  sync.Pool
	reset func(T) T
}

// New is equal to NewPool(new, nil).
func New[T any](new func() T) *Pool[T] {
	return NewPool(new, nil)
}

// NewPool returns a new pool.
//
// If reset is nil, do nothing when putting the object back into the pool.
func NewPool[T any](new func() T, reset func(T) T) *Pool[T] {
	pool := &Pool[T]{reset: reset}
	pool.pool.New = func() interface{} { return &Object[T]{pool: pool, Object: new()} }
	return pool
}

// Get returns an object from the pool.
func (p *Pool[T]) Get() *Object[T] {
	return p.pool.Get().(*Object[T])
}

// Put puts the object back into the pool.
func (p *Pool[T]) Put(o *Object[T]) {
	if o == nil {
		return
	}

	if o.pool != p {
		panic("the object is not allocated from the pool")
	}

	if p.reset != nil {
		o.Object = p.reset(o.Object)
	}

	p.pool.Put(o)
}
