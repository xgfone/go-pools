// Copyright 2020 xgfone
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
)

var (
	// BytesPool1K the bytes pool with 1K buffer.
	BytesPool1K = NewBytesPool(1024)

	// BytesPool2K the bytes pool with 2K buffer.
	BytesPool2K = NewBytesPool(2048)

	// BytesPool4K the bytes pool with 4K buffer.
	BytesPool4K = NewBytesPool(4096)

	// BytesPool8K the bytes pool with 8K buffer.
	BytesPool8K = NewBytesPool(8192)
)

// BytesPool is the []byte wrapper of sync.Pool.
type BytesPool struct {
	pool sync.Pool
}

// NewBytesPool returns a new []byte pool.
//
// size is the size of the []byte.
func NewBytesPool(size int) *BytesPool {
	newf := func() interface{} { return make([]byte, size) }
	return &BytesPool{pool: sync.Pool{New: newf}}
}

// Get returns a []byte.
func (p *BytesPool) Get() []byte {
	return p.pool.Get().([]byte)
}

// Put places a []byte to the pool.
func (p *BytesPool) Put(b []byte) {
	if len(b) != 0 {
		p.pool.Put(b)
	}
}
