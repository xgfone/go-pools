// Copyright 2022~2023 xgfone
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
	if new == nil {
		panic("NewPool: the new function must not be nil")
	}

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

// CapPool is an object pool to allocate the object based on the cap,
// such as a buffer or slice, from a befitting pool.
type CapPool[T any] struct {
	pool8   *Pool[T]
	pool16  *Pool[T]
	pool32  *Pool[T]
	pool64  *Pool[T]
	pool128 *Pool[T]
	pool256 *Pool[T]
	pool512 *Pool[T]
	pool1K  *Pool[T]
	pool2K  *Pool[T]
	pool4K  *Pool[T]
	pool8K  *Pool[T]
	pool16K *Pool[T]
	pool32K *Pool[T]
}

// NewCapPool returns a new CapPool.
//
// new is mandatory and reset is optional.
func NewCapPool[T any](new func(cap int) T, reset func(T) T) *CapPool[T] {
	if new == nil {
		panic("NewCapPool: the new function must not be nil")
	}

	return &CapPool[T]{
		pool8:   NewPool(func() T { return new(8) }, reset),
		pool16:  NewPool(func() T { return new(16) }, reset),
		pool32:  NewPool(func() T { return new(32) }, reset),
		pool64:  NewPool(func() T { return new(64) }, reset),
		pool128: NewPool(func() T { return new(128) }, reset),
		pool256: NewPool(func() T { return new(256) }, reset),
		pool512: NewPool(func() T { return new(512) }, reset),
		pool1K:  NewPool(func() T { return new(1024) }, reset),
		pool2K:  NewPool(func() T { return new(2048) }, reset),
		pool4K:  NewPool(func() T { return new(4096) }, reset),
		pool8K:  NewPool(func() T { return new(8192) }, reset),
		pool16K: NewPool(func() T { return new(16384) }, reset),
		pool32K: NewPool(func() T { return new(32768) }, reset),
	}
}

// Get returns an object from the befitting pool, which can be released
// into the original pool by calling the Release method of the returned object.
func (p *CapPool[T]) Get(cap int) *Object[T] {
	switch {
	case cap <= 8:
		return p.pool8.Get()

	case cap <= 16:
		return p.pool16.Get()

	case cap <= 32:
		return p.pool32.Get()

	case cap <= 64:
		return p.pool64.Get()

	case cap <= 128:
		return p.pool128.Get()

	case cap <= 256:
		return p.pool256.Get()

	case cap <= 1024:
		return p.pool1K.Get()

	case cap <= 2048:
		return p.pool2K.Get()

	case cap <= 4096:
		return p.pool4K.Get()

	case cap <= 8192:
		return p.pool8K.Get()

	case cap <= 16384:
		return p.pool16K.Get()

	default:
		return p.pool32K.Get()
	}
}
