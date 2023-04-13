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

import (
	"bytes"
	"fmt"
	"testing"
)

func ExamplePool() {
	type Context struct {
		// ....
	}
	pool := New(func() *Context { return new(Context) })

	// Get the context from the pool.
	ctx := pool.Get()

	// Use the object as *Context to do something.
	fmt.Println(ctx.Object) // ctx.Object => *Context
	// ...

	// Release the context into the pool.
	ctx.Release()

	// Output:
	// &{}
}

func TestCapPool(t *testing.T) {
	// For *bytes.Buffer
	bufferPool := NewCapPool(func(cap int) *bytes.Buffer {
		return bytes.NewBuffer(make([]byte, 0, cap))
	}, func(buf *bytes.Buffer) *bytes.Buffer {
		buf.Reset()
		return buf
	})
	if cap := bufferPool.Get(4).Object.Cap(); cap != 8 {
		t.Errorf("expect cap %d, but got %d", 8, cap)
	}
	if cap := bufferPool.Get(8).Object.Cap(); cap != 8 {
		t.Errorf("expect cap %d, but got %d", 8, cap)
	}
	if cap := bufferPool.Get(10).Object.Cap(); cap != 16 {
		t.Errorf("expect cap %d, but got %d", 16, cap)
	}

	// For []byte or []interface{}
	slicePool := NewCapPool(
		func(cap int) []byte { return make([]byte, 0, cap) },
		func(buf []byte) []byte { return buf[:0] },
	)
	if cap := cap(slicePool.Get(4).Object); cap != 8 {
		t.Errorf("expect cap %d, but got %d", 8, cap)
	}
	if cap := cap(slicePool.Get(8).Object); cap != 8 {
		t.Errorf("expect cap %d, but got %d", 8, cap)
	}
	if cap := cap(slicePool.Get(10).Object); cap != 16 {
		t.Errorf("expect cap %d, but got %d", 16, cap)
	}
}
