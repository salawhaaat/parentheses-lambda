// main.go
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/salawhaaat/parentheses-lambda/pkg/parentheses"
)

// Base URL for the API
const baseURL = "https://y3op6n0u1i.execute-api.eu-north-1.amazonaws.com/parentheses-byhand"

func main() {
	numberOfRequests := 1000
	lengths := []int{2, 4, 8}
	numWorkers := 10

	for _, length := range lengths {
		result := make(chan int, numberOfRequests) // Make result channel with buffer size
		jobs := make(chan string, numberOfRequests)

		// Create the workers
		var wg sync.WaitGroup
		for i := 0; i < numWorkers; i++ {
			wg.Add(1)
			go work(jobs, result, &wg)
		}

		url := fmt.Sprintf("%s?n=%d", baseURL, length)

		for i := 0; i < numberOfRequests; i++ {
			jobs <- url // Send the jobs
		}
		close(jobs)

		go func() {
			wg.Wait()
			close(result)
		}()

		successCount := 0

		for success := range result {
			successCount += success // Gather the results
		}

		successRate := float64(successCount) / float64(numberOfRequests) * 100
		fmt.Printf("Result for %d requests with length %d is %.2f%% \n", numberOfRequests, length, successRate)
	}
}

func work(jobs <-chan string, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		sequence, err := sendRequest(job)
		if err != nil {
			fmt.Println("Request error:", err)
			continue
		}
		if parentheses.IsBalanced(sequence) {
			results <- 1
		}
	}
}


func sendRequest(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response struct {
		Sequence string `json:"sequence"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("error parsing JSON response: %w", err)
	}

	return response.Sequence, nil
}
