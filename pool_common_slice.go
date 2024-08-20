// Copyright 2024 xgfone
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

var (
	// InterfacesPool is the pre-defined []any pool.
	InterfacesPool = NewCapSlicePool(
		func(cap int) []any { return make([]any, 0, cap) },
		func(vs []any) []any { return vs[:0] },
	)
)

// GetInterfaces returns an interfaces with len==0 from the befitting pool,
// which can be released into the original pool by calling the release function.
func GetInterfaces(cap int) *SliceObject[[]any] {
	return InterfacesPool.Get(cap)
}
