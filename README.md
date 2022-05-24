# workqueue

Package workqueue provides a means to queueing work.

## Documentation

[GoDoc](https://pkg.go.dev/github.com/nesv/workqueue/v2)

## Example

```go
package main

import (
	"fmt"
	"sync"

	workqueue "github.com/nesv/workqueue/v2"
)

func main() {
	// Create a new Queue.
	q := workqueue.New(1024)

	// Use a sync.WaitGroup to make sure we process all work before
	// exiting.
	var wg sync.WaitGroup

	// Now, let's do some work.
	for i := 0; i < 2048; i++ {
		wg.Add(1)
		go func(v int) {
			q <- func() {
				fmt.Println(v)
				wg.Done()
			}
		}(i)
	}

	// Wait for all the work to be done, then close the Queue.
	wg.Wait()
	close(q)
}
```
