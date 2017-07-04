package workqueue

// Copyright 2017 Nick Saika
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func TestWorkQueue(t *testing.T) {
	var wg sync.WaitGroup
	for i := 10; i < 1000000; i *= 10 {
		wg.Add(i)
		q := New(i)
		for j := 0; j < i; j++ {
			go func(w int) {
				q <- func() {
					dur := time.Duration(rand.Intn(10))
					time.Sleep(dur * time.Millisecond)
					t.Log(w)
					wg.Done()
				}
			}(j)
		}
		wg.Wait()
		close(q)
	}
}

func ExampleNew() {
	// Create a new WorkQueue.
	wq := New(1024)

	// This sync.WaitGroup is to make sure we wait until all of our work
	// is done.
	var wg sync.WaitGroup

	// Do some work.
	for i := 0; i < 2048; i++ {
		wg.Add(1)
		go func(v int) {
			wq <- func() {
				defer wg.Done()

				time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
				fmt.Println(v)
			}
		}(i)
	}

	// Wait for all of the work to finish, then close the WorkQueue.
	wg.Wait()
	close(wq)
}
