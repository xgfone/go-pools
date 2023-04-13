# go-pools [![Build Status](https://github.com/xgfone/go-pools/actions/workflows/go.yml/badge.svg)](https://github.com/xgfone/go-pools/actions/workflows/go.yml) [![GoDoc](https://pkg.go.dev/badge/github.com/xgfone/go-pools)](https://pkg.go.dev/github.com/xgfone/go-pools) [![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg?style=flat-square)](https://raw.githubusercontent.com/xgfone/go-pools/master/LICENSE)

Provide an object pool based on the generics supporting `Go1.18+`, such as `Pool`, `CapPool`.


## Install
```shell
$ go get -u github.com/xgfone/go-pools
```


## Example

#### `Pool`
```go
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
```

#### `CapPool`
```go
// For *bytes.Buffer
bufferPool := NewCapPool(
    func(cap int) *bytes.Buffer { return bytes.NewBuffer(make([]byte, 0, cap)) }, // new
    func(buf *bytes.Buffer) *bytes.Buffer { buf.Reset(); return buf },            // reset
)

buffer := bufferPool.Get(1024)
// TODO ...
buffer.Release()

// For a slice, such as []byte or []interface{}
slicePool := NewCapPool(
    func(cap int) []byte { return make([]byte, 0, cap) }, // new
    func(buf []byte) []byte { return buf[:0] },           // reset
)

bytes := slicePool.Get(128)
// TODO ...
bytes.Release()
```
