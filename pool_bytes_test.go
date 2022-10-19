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

func BenchmarkBytesPool(b *testing.B) {
	pool := NewBytesPool(8)
	pool.Get().Release()

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			pool.Get().Release()
		}
	})
}

func ExampleNewBytesPool() {
	pool := NewBytesPool(8)

	// Get the []byte object.
	bytes := pool.Get()

	// Use []byte to do something.
	fmt.Println(bytes.Object) // bytes.Object => []byte

	// Release the []byte object into the pool.
	bytes.Release()

	// Output:
	// []
}

func TestFixedBytesPool(t *testing.T) {
	bytes := FixedBytesPool64.Get()
	if len(bytes.Object) != 64 {
		t.Errorf("expect %d size, but got %d", 64, len(bytes.Object))
	}

	bytes.Release()
	bytes = FixedBytesPool64.Get()
	if len(bytes.Object) != 64 {
		t.Errorf("expect %d size, but got %d", 64, len(bytes.Object))
	}
}
