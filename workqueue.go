// Package workqueue provides an implementation of a queue for performing
// tasks with a number of background worker processes. At its core, this
// package utilizes a lot of the inherent properties of channels.
package workqueue

// WorkQueue is a channel type that you can send Work on.
type WorkQueue chan Work

// New creates and returns a new WorkQueue.
func New(numWorkers int) WorkQueue {
	queue := make(WorkQueue)
	d := make(dispatcher, numWorkers)
	go d.dispatch(queue)
	return queue
}

// Work is a task to perform that can be sent over a WorkQueue.
type Work func()

type dispatcher chan chan Work

func newDispatcher(queue WorkQueue, numWorkers int) dispatcher {
	d := make(dispatcher, numWorkers)
	go d.dispatch(queue)
	return d
}

func (d dispatcher) dispatch(queue WorkQueue) {
	// Create and start all of our workers.
	for i := 0; i < cap(d); i++ {
		w := make(worker)
		go w.work(d)
	}

	// Start the main loop in a goroutine.
	go func() {
		for work := range queue {
			go func(work Work) {
				worker := <-d
				worker <- work
			}(work)
		}

		// If we get here, the work queue has been closed, and we should
		// stop all of the workers.
		for i := 0; i < cap(d); i++ {
			w := <-d
			close(w)
		}
	}()
}

type worker chan Work

func (w worker) work(d dispatcher) {
	// Add ourselves to the dispatcher.
	d <- w

	// Start the main loop.
	go w.wait(d)
}

func (w worker) wait(d dispatcher) {
	for work := range w {
		// Do the work.
		if work == nil {
			panic("nil work received")
		}

		work()

		// Re-add ourselves to the dispatcher.
		d <- w
	}
}
