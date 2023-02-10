
# Workerpool
[![MIT License](https://img.shields.io/badge/License-MIT-green.svg)](https://choosealicense.com/licenses/mit/)

Simple constructor of workerpools

## Installation

Install the package with:

```bash
go get github.com/shtorm-7/workerpool
```
    
## Examples

Simple example usage of a workerpool:
```go
package main

import (
	"fmt"

	"github.com/shtorm-7/workerpool/constant"
	"github.com/shtorm-7/workerpool/pool"
	"github.com/shtorm-7/workerpool/worker"
)

func main() {
	queue := make(constant.Queue, 10)
	wp := pool.NewResizablePool(
		worker.NewWorkerFactory(
			queue,
			worker.WithFlow(worker.GracefulFlow),
		),
		10,
	)
	wp.Start()
	for i := 0; i < 10; i++ {
		i := i
		queue <- func() {
			fmt.Println(i)
		}
	}
	wp.Stop()
}
```

Also, you can use ```tools``` package for easier task creation

If you need to await your task, you can use ```tools.Await```

```go
package main

import (
	"github.com/shtorm-7/workerpool/constant"
	"github.com/shtorm-7/workerpool/pool"
	"github.com/shtorm-7/workerpool/tools"
	"github.com/shtorm-7/workerpool/worker"
)

func main() {
	queue := make(constant.Queue, 10)
	wp := pool.NewResizablePool(
		worker.NewWorkerFactory(
			queue,
			worker.WithFlow(worker.GracefulFlow),
		),
		10,
	)
	wp.Start()
	<-tools.Await(
		queue,
		func() {
			// your code
		},
	)
	wp.Stop()
}
```

Or ```tools.Future```, if you need to get a result from the task:

```go
package main

import (
	"fmt"

	"github.com/shtorm-7/workerpool/constant"
	"github.com/shtorm-7/workerpool/pool"
	"github.com/shtorm-7/workerpool/tools"
	"github.com/shtorm-7/workerpool/worker"
)

func main() {
	queue := make(constant.Queue, 10)
	wp := pool.NewResizablePool(
		worker.NewWorkerFactory(
			queue,
			worker.WithFlow(worker.GracefulFlow),
		),
		10,
	)
	wp.Start()
	taskResult := <-tools.Future(
		queue,
		func() (string, error) {
			// your code
			return "some result", nil
		},
	)
	fmt.Println(taskResult.Result, taskResult.Err)
	wp.Stop()
}
```

To execute multiple tasks, you can use ```tools.Chain``` (Or ```tools.Batch``` if you have custom tasks):

```go
package main

import (
	"fmt"

	"github.com/shtorm-7/workerpool/constant"
	"github.com/shtorm-7/workerpool/generator"
	"github.com/shtorm-7/workerpool/pool"
	"github.com/shtorm-7/workerpool/tools"
	"github.com/shtorm-7/workerpool/worker"
)

func main() {
	queue := make(constant.Queue, 10)
	wp := pool.NewResizablePool(
		worker.NewWorkerFactory(
			queue,
			worker.WithFlow(worker.GracefulFlow),
		),
		10,
	)
	wp.Start()
	chain := tools.NewChain(
		tools.NewLink[string, string](
			queue,
			func(value string) (string, error) {
				// your code
				return value, nil
			},
		),
	)
	for taskResult := range chain.Batch(
		10, generator.Range([]string{"first value", "second value"}),
	) {
		fmt.Println(taskResult.Result)
	}
	wp.Stop()
}
```

if you need to separate your logic into different workerpools, then you can add more  ```tools.Link``` to the ```tools.Chain```:

```go
package main

import (
	"fmt"

	"github.com/shtorm-7/workerpool/constant"
	"github.com/shtorm-7/workerpool/generator"
	"github.com/shtorm-7/workerpool/pool"
	"github.com/shtorm-7/workerpool/tools"
	"github.com/shtorm-7/workerpool/worker"
)

func main() {
	queue1 := make(constant.Queue, 10)
	queue2 := make(constant.Queue, 10)
	wp := pool.NewPool(
		[]constant.WorkerFactory{
			pool.NewResizablePoolFactory(
				worker.NewWorkerFactory(
					queue1,
					worker.WithFlow(worker.GracefulFlow),
				),
				10,
			),
			pool.NewResizablePoolFactory(
				worker.NewWorkerFactory(
					queue2,
					worker.WithFlow(worker.GracefulFlow),
				),
				10,
			),
		},
	)
	wp.Start()
	chain := tools.NewChain(
		tools.AddLink(
			tools.NewLink[int, int](queue1,
				func(value int) (string, error) {
					// your code
					return "some result", nil
				},
			), queue2,
			func(value string) (int, error) {
				// your code
				return 0, nil
			},
		),
	)
	for taskResult := range chain.Batch(
		100, generator.SequenceRange(0, 10),
	) {
		fmt.Println(taskResult.Result)
	}
	wp.Stop()
}
```

## TODO

* Publish stable release
* Add tests
* Add logger
* Add more documentation

## License

This package is released under the MIT license. See the complete license in the package
