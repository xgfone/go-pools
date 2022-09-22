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

// Pre-define some fixed bytes pools with the different size.
var (
	FixedBytesPool64  = NewFixedBytesPool(64)
	FixedBytesPool128 = NewFixedBytesPool(128)
	FixedBytesPool256 = NewFixedBytesPool(256)
	FixedBytesPool512 = NewFixedBytesPool(512)
	FixedBytesPool1K  = NewFixedBytesPool(1024)
	FixedBytesPool2K  = NewFixedBytesPool(2048)
	FixedBytesPool4K  = NewFixedBytesPool(4096)
	FixedBytesPool8K  = NewFixedBytesPool(8192)
)

// GetFixedBytes returns a fixed bytes from the befitting pool, which can
// be released into the original pool by calling the release function.
func GetFixedBytes(size int) *FixedBytes {
	if size <= 64 {
		return FixedBytesPool64.Get()
	} else if size <= 128 {
		return FixedBytesPool128.Get()
	} else if size <= 256 {
		return FixedBytesPool256.Get()
	} else if size <= 512 {
		return FixedBytesPool512.Get()
	} else if size <= 1024 {
		return FixedBytesPool1K.Get()
	} else if size <= 2048 {
		return FixedBytesPool2K.Get()
	} else if size <= 4096 {
		return FixedBytesPool4K.Get()
	} else {
		return FixedBytesPool8K.Get()
	}
}

// FixedBytes is used to enclose the fixed byte slice []byte.
type FixedBytes struct {
	Bytes []byte
	pool  *FixedBytesPool
}

// Release releases the fixed bytes into the original pool.
func (b *FixedBytes) Release() {
	if b != nil && b.pool != nil {
		b.pool.Put(b)
	}
}

// FixedBytesPool is the pool to allocate the fixed bytes.
type FixedBytesPool struct{ pool sync.Pool }

// NewFixedBytesPool returns a new fixed bytes pool.
func NewFixedBytesPool(size int) *FixedBytesPool {
	pool := new(FixedBytesPool)
	pool.pool.New = func() interface{} {
		return &FixedBytes{pool: pool, Bytes: make([]byte, size)}
	}
	return pool
}

// Get returns a fixed bytes from the pool.
func (p *FixedBytesPool) Get() *FixedBytes { return p.pool.Get().(*FixedBytes) }

// Put puts the fixed bytes back into the pool.
func (p *FixedBytesPool) Put(b *FixedBytes) { p.pool.Put(b) }
