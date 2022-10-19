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

	FixedBytesPool64  = NewFixedBytesPool(64)
	FixedBytesPool128 = NewFixedBytesPool(128)
	FixedBytesPool256 = NewFixedBytesPool(256)
	FixedBytesPool512 = NewFixedBytesPool(512)
	FixedBytesPool1K  = NewFixedBytesPool(1024)
	FixedBytesPool2K  = NewFixedBytesPool(2048)
	FixedBytesPool4K  = NewFixedBytesPool(4096)
	FixedBytesPool8K  = NewFixedBytesPool(8192)
)

// NewBytesPool returns a new pool based on []byte.
func NewBytesPool(cap int) *Pool[[]byte] {
	return NewPool(func() []byte {
		return make([]byte, 0, cap)
	}, func(b []byte) []byte {
		return b[:0]
	})
}

// NewFixedBytesPool returns a new pool based on the fixed-size []byte.
func NewFixedBytesPool(size int) *Pool[[]byte] {
	return New(func() []byte {
		return make([]byte, size)
	})
}

// GetBytes returns a bytes from the befitting pool, which can be released
// into the original pool by calling the release function.
func GetBytes(cap int) *Object[[]byte] {
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

// GetFixedBytes returns a fixed bytes from the befitting pool, which can
// be released into the original pool by calling the release function.
func GetFixedBytes(size int) *Object[[]byte] {
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
