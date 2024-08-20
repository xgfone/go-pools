// Copyright 2024 xgfone
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

// SliceObject is used to enclose T.
type SliceObject[T ~[]any] struct {
	Objects T

	pool *SlicePool[T]
}

// Get returns the field Objects.
func (o *SliceObject[T]) Get() T { return o.Objects }

// Len returns the length of the slice objects.
func (o *SliceObject[T]) Len() int {
	if o == nil {
		return 0
	}
	return len(o.Objects)
}

// Append appends some objects into the slice objects and return itself.
func (o *SliceObject[T]) Append(objects ...any) *SliceObject[T] {
	o.Objects = append(o.Objects, objects...)
	return o
}

// Release releases the slice object into the original pool.
func (o *SliceObject[T]) Release() {
	if o != nil && o.pool != nil {
		o.pool.Put(o)
	}
}

// SlicePool is the slice object pool to allocate a slice object.
type SlicePool[T ~[]any] struct {
	pool  sync.Pool
	reset func(T) T
}

// NewSlicePool returns a new slice object pool.
//
// If reset is nil, do nothing when putting the slice object back into the pool.
func NewSlicePool[T ~[]any](new func() T, reset func(T) T) *SlicePool[T] {
	if new == nil {
		panic("NewSlicePool: the new function must not be nil")
	}

	pool := &SlicePool[T]{reset: reset}
	pool.pool.New = func() any { return &SliceObject[T]{pool: pool, Objects: new()} }
	return pool
}

// Get returns a slice object from the pool.
func (p *SlicePool[T]) Get() *SliceObject[T] {
	return p.pool.Get().(*SliceObject[T])
}

// Put puts the slice object back into the pool.
func (p *SlicePool[T]) Put(o *SliceObject[T]) {
	if o == nil {
		return
	}

	if o.pool != p {
		panic("the slice object is not allocated from the pool")
	}

	if p.reset != nil {
		o.Objects = p.reset(o.Objects)
	}

	p.pool.Put(o)
}

// CapSlicePool is a slice object pool to allocate the slice object
// based on the cap.
type CapSlicePool[T ~[]any] struct {
	pool8   *SlicePool[T]
	pool16  *SlicePool[T]
	pool32  *SlicePool[T]
	pool64  *SlicePool[T]
	pool128 *SlicePool[T]
	pool256 *SlicePool[T]
	pool512 *SlicePool[T]
	pool1K  *SlicePool[T]
	pool2K  *SlicePool[T]
	pool4K  *SlicePool[T]
	pool8K  *SlicePool[T]
	pool16K *SlicePool[T]
	pool32K *SlicePool[T]
}

// NewCapSlicePool returns a new CapSlicePool.
//
// new is mandatory and reset is optional.
func NewCapSlicePool[T ~[]any](new func(cap int) T, reset func(T) T) *CapSlicePool[T] {
	if new == nil {
		panic("NewCapSlicePool: the new function must not be nil")
	}

	return &CapSlicePool[T]{
		pool8:   NewSlicePool(func() T { return new(8) }, reset),
		pool16:  NewSlicePool(func() T { return new(16) }, reset),
		pool32:  NewSlicePool(func() T { return new(32) }, reset),
		pool64:  NewSlicePool(func() T { return new(64) }, reset),
		pool128: NewSlicePool(func() T { return new(128) }, reset),
		pool256: NewSlicePool(func() T { return new(256) }, reset),
		pool512: NewSlicePool(func() T { return new(512) }, reset),
		pool1K:  NewSlicePool(func() T { return new(1024) }, reset),
		pool2K:  NewSlicePool(func() T { return new(2048) }, reset),
		pool4K:  NewSlicePool(func() T { return new(4096) }, reset),
		pool8K:  NewSlicePool(func() T { return new(8192) }, reset),
		pool16K: NewSlicePool(func() T { return new(16384) }, reset),
		pool32K: NewSlicePool(func() T { return new(32768) }, reset),
	}
}

// Get returns a slice object from the befitting pool, which can be released
// into the original pool by calling the Release method of the returned object.
func (p *CapSlicePool[T]) Get(cap int) *SliceObject[T] {
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
