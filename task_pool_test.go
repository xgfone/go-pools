// Copyright 2019 xgfone
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
	"context"
	"fmt"
	"time"
)

func ExampleTaskPool() {
	task := func(ctx context.Context, args ...interface{}) (interface{}, error) {
		time.Sleep(args[0].(time.Duration))
		return args[0], nil
	}

	pool := NewTaskPool(3)

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
	// WaitAllTaskResults(r1, r2, r3, r4, r5, r6)

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
