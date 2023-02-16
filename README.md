
# Workerpool
[![Go Reference](https://pkg.go.dev/badge/github.com/shtorm-7/workerpool.svg)](https://pkg.go.dev/github.com/shtorm-7/workerpool)
[![Go Report Card](https://goreportcard.com/badge/github.com/shtorm-7/workerpool)](https://goreportcard.com/report/github.com/shtorm-7/workerpool)
[![MIT License](https://img.shields.io/badge/License-MIT-green.svg)](https://github.com/shtorm-7/workerpool/blob/main/LICENSE)

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
		worker.NewWorkerFactory(queue),
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
		worker.NewWorkerFactory(queue),
		10,
	)
	wp.Start()
	result := <-tools.Future(
		queue,
		func() string {
			// your code
			return "some result"
		},
	)
	fmt.Println(result)
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
		worker.NewWorkerFactory(queue),
		10,
	)
	wp.Start()
	chain := tools.NewChain(
		tools.NewLink[string](
			queue,
			func(value string) (string, error) {
				// your code
				return value, nil
			},
		),
	)
	for chainResult := range chain.Batch(
		10, generator.Range([]string{"first value", "second value"}),
	) {
		fmt.Println(chainResult.Result)
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
				worker.NewWorkerFactory(queue1),
				10,
			),
			pool.NewResizablePoolFactory(
				worker.NewWorkerFactory(queue2),
				10,
			),
		},
	)
	wp.Start()
	chain := tools.NewChain(
		tools.AddLink(
			tools.NewLink[string](queue1,
				func(value int) (int, error) {
					// your code
					return value, nil
				},
			), queue2,
			func(value int) (string, error) {
				// your code
				return fmt.Sprintf("result: %d", value), nil
			},
		),
	)
	for chainResult := range chain.Batch(
		10, generator.SequenceRange(0, 10),
	) {
		fmt.Println(chainResult.Result)
	}
	wp.Stop()
}
```

## TODO

* Publish stable release
* Add tests
* Add more documentation

## License

This package is released under the MIT license. See the complete license in the package
