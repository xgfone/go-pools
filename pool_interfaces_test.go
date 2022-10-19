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

func BenchmarkInterfacesPool(b *testing.B) {
	pool := NewInterfacesPool(8)
	pool.Get().Release()

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			pool.Get().Release()
		}
	})
}

func ExampleNewInterfacesPool() {
	pool := NewInterfacesPool(8)

	// Get the []interface{} object.
	interfaces := pool.Get()

	// Use []interface{} to do something.
	fmt.Println(interfaces.Object) // interfaces.Object => []interface{}

	// Release the []interface{} object into the pool.
	interfaces.Release()

	// Output:
	// []
}
