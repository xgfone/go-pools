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

import "fmt"

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
