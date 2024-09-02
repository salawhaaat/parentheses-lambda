package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/salawhaaat/parentheses-lambda/pkg/parentheses"
)

// API endpoint URL
const apiURL = "https://y3op6n0u1i.execute-api.eu-north-1.amazonaws.com/parentheses-byhand"

type Response struct {
	Sequence string `json:"sequence"`
}

func main() {
	numberOfRequests := 1000
	lengths := []int{2, 4, 8}
	concurrency := 7 // Number of concurrent workers

	for _, length := range lengths {
		startTime := time.Now()
		success := evaluateParenthesesWithWorkerPool(length, numberOfRequests, concurrency)
		duration := time.Since(startTime)

		successRate := float64(success) / float64(numberOfRequests) * 100
		fmt.Printf("Result for %d requests with length %d is %.2f%% (Time taken: %s)\n", numberOfRequests, length, successRate, duration)
	}
}

func evaluateParenthesesWithWorkerPool(length, numberOfRequests, concurrency int) int {
	var successCount int
	var wg sync.WaitGroup
	var mu sync.Mutex
	url := fmt.Sprintf("%s?n=%d", apiURL, length)

	// Channel to limit concurrency
	requests := make(chan struct{}, concurrency)

	for i := 0; i < numberOfRequests; i++ {
		requests <- struct{}{}
		wg.Add(1)

		go func() {
			defer wg.Done()

			sequence, err := sendRequest(url)
			if err != nil {
				fmt.Println("Request error:", err)
				<-requests

				return
			}

			if parentheses.IsBalanced(sequence) {
				mu.Lock()
				successCount++
				mu.Unlock()
			}

			<-requests
		}()
	}

	wg.Wait()

	return successCount
}

func sendRequest(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("Error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", fmt.Errorf("error parsing JSON response: %w", err)
	}

	return response.Sequence, nil
}
