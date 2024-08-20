// Copyright 2023 xgfone
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

// Pre-define some object-capacity pools, such as []byte, []any or *bytes.Buffer.
var (
	BufferPool = NewCapPool(
		func(cap int) *bytes.Buffer { return bytes.NewBuffer(make([]byte, 0, cap)) },
		func(buf *bytes.Buffer) *bytes.Buffer { buf.Reset(); return buf },
	)

	InterfacesPool = NewCapPool(
		func(cap int) []any { return make([]any, 0, cap) },
		func(vs []any) []any { return vs[:0] },
	)

	BytesPool = NewCapPool(
		func(cap int) []byte { return make([]byte, 0, cap) },
		func(bs []byte) []byte { return bs[:0] },
	)

	LenBytesPool = NewCapPool(
		func(cap int) []byte { return make([]byte, cap) },
		nil, // Use the original byte slice, which is equal to func(bs []byte) []byte { return bs }
	)
)

// GetBuffer returns a buffer from the befitting pool, which can be released
// into the original pool by calling the release function.
func GetBuffer(cap int) *Object[*bytes.Buffer] {
	return BufferPool.Get(cap)
}

// GetInterfaces returns an interfaces with len==0 from the befitting pool,
// which can be released into the original pool by calling the release function.
func GetInterfaces(cap int) *Object[[]any] {
	return InterfacesPool.Get(cap)
}

// GetBytes returns a bytes with len==0 from the befitting pool,
// which can be released into the original pool by calling the release function.
func GetBytes(cap int) *Object[[]byte] {
	return BytesPool.Get(cap)
}

// GetLenBytes returns a bytes with len==cap from the befitting pool,
// which can be released into the original pool by calling the release function.
func GetLenBytes(cap int) *Object[[]byte] {
	return LenBytesPool.Get(cap)
}
