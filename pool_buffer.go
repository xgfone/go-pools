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

import "bytes"

// Pre-define some buffer pools with the different capacity.
var (
	BufferPool64  = NewBufferPool(64)
	BufferPool128 = NewBufferPool(128)
	BufferPool256 = NewBufferPool(256)
	BufferPool512 = NewBufferPool(512)
	BufferPool1K  = NewBufferPool(1024)
	BufferPool2K  = NewBufferPool(2048)
	BufferPool4K  = NewBufferPool(4096)
	BufferPool8K  = NewBufferPool(8192)
)

// NewBufferPool returns a new pool based on *bytes.Buffer.
func NewBufferPool(cap int) *Pool[*bytes.Buffer] {
	return NewPool(func() *bytes.Buffer {
		return bytes.NewBuffer(make([]byte, 0, cap))
	}, func(b *bytes.Buffer) *bytes.Buffer {
		b.Reset()
		return b
	})
}

// GetBuffer returns a buffer from the befitting pool, which can be released
// into the original pool by calling the release function.
func GetBuffer(cap int) *Object[*bytes.Buffer] {
	if cap <= 64 {
		return BufferPool64.Get()
	} else if cap <= 128 {
		return BufferPool128.Get()
	} else if cap <= 256 {
		return BufferPool256.Get()
	} else if cap <= 512 {
		return BufferPool512.Get()
	} else if cap <= 1024 {
		return BufferPool1K.Get()
	} else if cap <= 2048 {
		return BufferPool2K.Get()
	} else if cap <= 4096 {
		return BufferPool4K.Get()
	} else {
		return BufferPool8K.Get()
	}
}
