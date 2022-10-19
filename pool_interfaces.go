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

// Pre-define some interfaces pools with the different capacity.
var (
	InterfacesPool8   = NewInterfacesPool(8)
	InterfacesPool16  = NewInterfacesPool(16)
	InterfacesPool32  = NewInterfacesPool(32)
	InterfacesPool64  = NewInterfacesPool(64)
	InterfacesPool128 = NewInterfacesPool(128)
	InterfacesPool256 = NewInterfacesPool(256)
	InterfacesPool512 = NewInterfacesPool(512)
	InterfacesPool1K  = NewInterfacesPool(1024)
)

// NewInterfacesPool returns a new pool based on []interface{}.
func NewInterfacesPool(cap int) *Pool[[]interface{}] {
	return NewPool(func() []interface{} {
		return make([]interface{}, 0, cap)
	}, func(b []interface{}) []interface{} {
		return b[:0]
	})
}

// GetInterfaces returns an interfaces from the befitting pool,
// which can be released into the original pool by calling the release function.
func GetInterfaces(cap int) *Object[[]interface{}] {
	if cap <= 8 {
		return InterfacesPool8.Get()
	} else if cap <= 16 {
		return InterfacesPool16.Get()
	} else if cap <= 32 {
		return InterfacesPool32.Get()
	} else if cap <= 64 {
		return InterfacesPool64.Get()
	} else if cap <= 128 {
		return InterfacesPool128.Get()
	} else if cap <= 256 {
		return InterfacesPool256.Get()
	} else if cap <= 512 {
		return InterfacesPool512.Get()
	} else {
		return InterfacesPool1K.Get()
	}
}
