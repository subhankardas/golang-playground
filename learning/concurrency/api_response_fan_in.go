package concurrency

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

// Write a Go program that launches three concurrent API calls and returns the result of whichever call completes first.
// As soon as the first response is received, the program should exit while ensuring the remaining goroutines are properly
// cancelled to prevent goroutine leaks.
func FanInAPIResponses() {
	ctx, cancel := context.WithCancel(context.Background()) // Cancelable context to manage API request lifetimes
	responses := make(chan string, 1)                       // Buffered channel to hold one response, prevents blocking

	defer cancel()

	go mockAPIRequestV2(ctx, 1, responses)
	go mockAPIRequestV2(ctx, 2, responses)
	go mockAPIRequestV2(ctx, 3, responses)

	res := <-responses // Collect the first response that arrives
	fmt.Println("Received:", res)

	cancel() // Cancel other API requests that are still running

	// Wait a bit to view cancellation logs
	time.Sleep(10 * time.Second)
	fmt.Println("Main exiting...")
}

func mockAPIRequestV2(ctx context.Context, id int, resCh chan string) {
	// Simulate variable response times
	delay := rand.Intn(3) + 1 // Random delay between 1 and 3 seconds

	select {
	case <-time.After(time.Duration(delay) * time.Second):
		select {
		case resCh <- fmt.Sprintf("Response from API %d", id): // Send response to channel

		case <-ctx.Done(): // Handle cancellation while sending
			fmt.Printf("API %d request cancelled while sending response\n", id)
			return
		}
	case <-ctx.Done(): // Handle cancellation while waiting
		fmt.Printf("API %d request cancelled\n", id)
		return
	}
}
