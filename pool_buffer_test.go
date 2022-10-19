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
	"fmt"
	"testing"
)

func BenchmarkBufferPool(b *testing.B) {
	pool := NewBufferPool(8)
	pool.Get().Release()

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			pool.Get().Release()
		}
	})
}

func ExampleNewBufferPool() {
	pool := NewBufferPool(8)

	// Get the *bytes.Buffer object.
	buf := pool.Get()

	// Use *bytes.Buffer to do something.
	fmt.Println(buf.Object) // buf.Object => *bytes.Buffer

	// Release the *bytes.Buffer object into the pool.
	buf.Release()

	// Output:
	//
}
