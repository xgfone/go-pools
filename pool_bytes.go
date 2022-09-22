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

// Pre-define some bytes pools with the different capacity.
var (
	BytesPool64  = NewBytesPool(64)
	BytesPool128 = NewBytesPool(128)
	BytesPool256 = NewBytesPool(256)
	BytesPool512 = NewBytesPool(512)
	BytesPool1K  = NewBytesPool(1024)
	BytesPool2K  = NewBytesPool(2048)
	BytesPool4K  = NewBytesPool(4096)
	BytesPool8K  = NewBytesPool(8192)
)

// GetBytes returns a bytes from the befitting pool, which can be released
// into the original pool by calling the release function.
func GetBytes(cap int) *Bytes {
	if cap <= 64 {
		return BytesPool64.Get()
	} else if cap <= 128 {
		return BytesPool128.Get()
	} else if cap <= 256 {
		return BytesPool256.Get()
	} else if cap <= 512 {
		return BytesPool512.Get()
	} else if cap <= 1024 {
		return BytesPool1K.Get()
	} else if cap <= 2048 {
		return BytesPool2K.Get()
	} else if cap <= 4096 {
		return BytesPool4K.Get()
	} else {
		return BytesPool8K.Get()
	}
}

// Bytes is used to enclose the byte slice []byte.
type Bytes struct {
	Bytes []byte
	pool  *BytesPool
}

// Release releases the bytes into the original pool.
func (b *Bytes) Release() {
	if b != nil && b.pool != nil {
		b.pool.Put(b)
	}
}

// BytesPool is the pool to allocate the bytes.
type BytesPool struct{ pool sync.Pool }

// NewBytesPool returns a new bytes pool.
func NewBytesPool(cap int) *BytesPool {
	pool := new(BytesPool)
	pool.pool.New = func() interface{} {
		return &Bytes{pool: pool, Bytes: make([]byte, 0, cap)}
	}
	return pool
}

// Get returns a bytes from the pool.
func (p *BytesPool) Get() *Bytes {
	b := p.pool.Get().(*Bytes)
	b.Bytes = b.Bytes[:0]
	return b
}

// Put puts the bytes back into the pool.
func (p *BytesPool) Put(b *Bytes) { p.pool.Put(b) }
