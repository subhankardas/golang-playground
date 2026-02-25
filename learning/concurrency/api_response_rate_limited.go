package concurrency

import (
	"fmt"
	"sync"
	"time"
)

// APIResponseRateLimited demonstrates how to implement a rate limiter for API responses using a
// ticker. This example simulates processing a number of tasks with a specified rate limit, ensuring
// that no more than the rate limit number of tasks are processed per second.
func APIResponseRateLimited() {
	var (
		numWorkers = 10 // no. of workers
		numJobs    = 20 // no. of jobs
		rateLimit  = 5  // rate limit in jobs per second
	)

	ticker := time.NewTicker(time.Second / time.Duration(rateLimit)) // ticks at the rate limit
	defer ticker.Stop()

	var tasks = make(chan string, numJobs)
	var wg sync.WaitGroup

	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go workerWithTicker(i, tasks, ticker.C, &wg)
	}

	for i := 1; i <= numJobs; i++ {
		tasks <- fmt.Sprintf("task: %d", i)
	}

	close(tasks)
	wg.Wait()
}

func workerWithTicker(id int, tasks <-chan string, limiter <-chan time.Time, wg *sync.WaitGroup) {
	defer wg.Done()

	for task := range tasks {
		<-limiter // wait for the next tick

		fmt.Printf("Worker: %d processing %s at %s\n", id, task, time.Now().Format("15:04:05.000"))
		time.Sleep(100 * time.Millisecond) // simulate API call
	}
}
