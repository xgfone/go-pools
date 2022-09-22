# go-pools [![Build Status](https://github.com/xgfone/go-pools/actions/workflows/go.yml/badge.svg)](https://github.com/xgfone/go-pools/actions/workflows/go.yml) [![GoDoc](https://pkg.go.dev/badge/github.com/xgfone/go-pools)](https://pkg.go.dev/github.com/xgfone/go-pools) [![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg?style=flat-square)](https://raw.githubusercontent.com/xgfone/go-pools/master/LICENSE)

This package supporting `Go1.7+` provides some pools, such as `BufferPool`, `BytesPool`, `FixedBytesPool`, `TaskPool`, etc.


## Install
```shell
$ go get -u github.com/xgfone/go-pools
```


## Example

```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/xgfone/go-pools"
)

func task(ctx context.Context, args ...interface{}) (interface{}, error) {
	time.Sleep(args[0].(time.Duration))
	return args[0], nil
}

func main() {
	pool := pools.NewTaskPool(3)

	// Terminate the task pool after some seconds.
	go func() {
		time.Sleep(time.Millisecond * 200)
		pool.Shutdown(nil) // Return immediately.
	}()

	// Run the tasks in the pool.
	r1 := pool.RunTaskFuncWithResult(task, time.Millisecond*10)
	r2 := pool.RunTaskFuncWithResult(task, time.Millisecond*20)
	r3 := pool.RunTaskFuncWithResult(task, time.Millisecond*30)
	r4 := pool.RunTaskFuncWithResult(task, time.Millisecond*40)
	r5 := pool.RunTaskFuncWithResult(task, time.Millisecond*50)
	r6 := pool.RunTaskFuncWithResult(task, time.Millisecond*60)

	pool.Wait() // Wait until the whole task pool exits.
	// Or, Wait until all the task terminate.
	// pools.WaitAllTaskResults(r1, r2, r3, r4, r5, r6)

	fmt.Println("task1 result:", r1.Result())
	fmt.Println("task2 result:", r2.Result())
	fmt.Println("task3 result:", r3.Result())
	fmt.Println("task4 result:", r4.Result())
	fmt.Println("task5 result:", r5.Result())
	fmt.Println("task6 result:", r6.Result())

	// Output:
	// task1 result: 10ms
	// task2 result: 20ms
	// task3 result: 30ms
	// task4 result: 40ms
	// task5 result: 50ms
	// task6 result: 60ms
}
```
