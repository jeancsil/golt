package main

import (
	"flag"
	"fmt"
	"net/http"
	"sort"
	"sync"
	"time"
)

func main() {
	// method := flag.String("m", "GET", "HTTP method to use")
	url := flag.String("u", "http://google.com", "URL to use")
	requests := flag.Int("n", 10, "Number of requests to make")
	concurrency := flag.Int("c", 1, "Number of concurrent requests")
	flag.Parse()

	var jobs = make(chan string)
	var results = make(chan time.Duration, *requests)

	var wg sync.WaitGroup
	wg.Add(*concurrency)
	for i := 0; i < *concurrency; i++ {
		go func() {
			for u := range jobs {
				results <- makeRequest(u)
			}
			wg.Done()
		}()
	}

	// send requests by sending url to jobs channel
	for i := 0; i < *requests; i++ {
		jobs <- *url
	}
	close(jobs)
	wg.Wait()

	// Close results channel so the range loop in the next step can finish
	close(results)

	timings := make([]time.Duration, 0, *requests)
	for r := range results {
		timings = append(timings, r)
	}

	sort.Slice(timings, func(i, j int) bool {
		return timings[i] < timings[j]
	})
	fmt.Printf("Timings: %v\n", timings)
	fmt.Printf("Fastest: %s\n", timings[0])
	fmt.Printf("Slowest: %s\n", timings[*requests-1])
	fmt.Printf("Median: %s\n", timings[*requests/2])

}

func makeRequest(url string) time.Duration {
	var start = time.Now()
	resp, err := http.Get(url)

	if err != nil {
		fmt.Println(err)
		return time.Duration(0)
	}

	fmt.Printf("Fetched %s with status: %s\n", url, resp.Status)
	defer resp.Body.Close()

	return time.Since(start)
}
