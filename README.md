# workqueue

Package workqueue provides a flexible means to queueing work.

## Documentation

[GoDoc](https://godoc.org/github.com/nesv/workqueue)

## Example

```go
package main

import (
	"fmt"
	"sync"

	"github.com/nesv/workqueue"
)

func main() {
	// Create a new WorkQueue.
	wq := New(1024)

	// Use a sync.WaitGroup to make sure we process all work before
	// exiting.
	var wg sync.WaitGroup

	// Now, let's do some work.
	for i := 0; i < 2048; i++ {
		wg.Add(1)
		go func(v int) {
			wq <- func() {
				fmt.Println(v)
				wg.Done()
			}
		}(i)
	}

	// Wait for all the work to be done, then close the WorkQueue.
	wg.Wait()
	close(wq)
}
```
