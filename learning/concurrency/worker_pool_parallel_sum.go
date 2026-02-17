package concurrency

import "sync"

func WorkerPoolParallelSum() {
	var n int = 100     // Sum of numbers from 1 to n
	var workers int = 4 // Number of worker goroutines

	var chunksize int = n / workers

	var dataCh chan int = make(chan int, workers) // To collect results from workers
	var wg sync.WaitGroup

	for i := 1; i <= workers; i++ {
		start := (chunksize * (i - 1)) + 1
		end := chunksize * i
		if i == workers {
			end = n
		}

		wg.Add(1)
		go worker(start, end, dataCh, &wg)
	}

	wg.Wait()
	close(dataCh)

	totalSum := 0
	for sum := range dataCh { // Collect partial sums from the channel
		totalSum += sum
	}
	println("Total Sum:", totalSum)
}
func worker(start, end int, dataCh chan int, wg *sync.WaitGroup) {
	defer wg.Done() // Mark this worker as done when the function returns

	sum := 0
	for i := start; i <= end; i++ {
		sum += i
	}

	dataCh <- sum // Send partial sum to the channel
}
