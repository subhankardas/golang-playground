package concurrency

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// You have 3 different functions, each mimicking an API call that returns a random integer after a random delay (1–3 seconds).
// The Goal: Write a function that starts all three "API calls" at once and returns a single channel. As soon as any API returns
// a value, it should be sent into that single channel. The main function should print them as they arrive.
// Tricky part: How do you close the destination channel only after all three APIs are done without blocking the main thread?
func AggregateAPIResponses() {
	var apiCount int = 3                     // Number of APIs to call
	responses := make(chan string, apiCount) // Buffered channel to hold API responses

	var wg sync.WaitGroup
	wg.Add(apiCount)

	for i := 1; i <= apiCount; i++ {
		go func(id int) {
			defer wg.Done()

			response := mockAPIRequest(id) // Simulate API request
			responses <- response          // Send response to channel
		}(i)
	}

	// Close in another goroutine, without blocking the main goroutine
	go func() {
		wg.Wait()
		close(responses)
	}()

	// Collect and print responses as they arrive
	for res := range responses {
		println(res)
	}
}

func mockAPIRequest(id int) string {
	// Simulate variable response times
	delay := rand.Intn(3) + 1 // Random delay between 1 and 3 seconds
	time.Sleep(time.Duration(delay) * time.Second)

	return fmt.Sprintf("Response from API %d", id)
}
