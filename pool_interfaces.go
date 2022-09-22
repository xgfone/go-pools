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

import "sync"

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

// GetInterfaces returns an interfaces from the befitting pool,
// which can be released into the original pool by calling the release function.
func GetInterfaces(cap int) *Interfaces {
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

// Interfaces is used to enclose []interface.
type Interfaces struct {
	Interfaces []interface{}
	pool       *InterfacesPool
}

// Release releases the interfaces into the original pool.
func (i *Interfaces) Release() {
	if i != nil && i.pool != nil {
		i.pool.Put(i)
	}
}

// InterfacesPool is the pool to allocate the interfaces.
type InterfacesPool struct{ pool sync.Pool }

// NewInterfacesPool returns a new interfaces pool.
func NewInterfacesPool(cap int) *InterfacesPool {
	pool := new(InterfacesPool)
	pool.pool.New = func() interface{} {
		return &Interfaces{pool: pool, Interfaces: make([]interface{}, 0, cap)}
	}
	return pool
}

// Get returns an interfaces from the pool.
func (p *InterfacesPool) Get() *Interfaces {
	i := p.pool.Get().(*Interfaces)
	i.Interfaces = i.Interfaces[:0]
	return i
}

// Put puts the interfaces back into the pool.
func (p *InterfacesPool) Put(i *Interfaces) { p.pool.Put(i) }
