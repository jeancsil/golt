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
	url := flag.String("u", "http://google.com", "URL to use")
	requests := flag.Int("n", 10, "Number of requests to make")
	concurrency := flag.Int("c", 1, "Number of concurrent requests")
	flag.Parse()

	var jobs = make(chan string)
	var results = make(chan time.Duration, *requests)

	var wg sync.WaitGroup
	wg.Add(*concurrency)
	var start = time.Now()
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
	timings = calculateStats(timings, results)
	reportStats(timings, start)
}

func reportStats(timings []time.Duration, start time.Time) {
	n_reqs := len(timings)
	percentiles := []int{50, 90, 95, 99}

	// fmt.Printf("Timings: %v\n", timings)
	fmt.Printf("\nFastest: %s\n", timings[0])
	fmt.Printf("Median: %s\n", timings[n_reqs/2])
	fmt.Printf("Slowest: %s\n", timings[n_reqs-1])
	fmt.Printf("Average: %f\n", float64(sliceSum(timings).Seconds())/float64(n_reqs))
	fmt.Printf("Throughput: %f req/s\n\n", float64(n_reqs)/time.Since(start).Seconds())

	fmt.Println("Percentiles:")
	length := len(timings)
	for _, p := range percentiles {
		index := int(float64(length) * float64(p) / 100)
		fmt.Printf("%.2f%%: %s\n", float64(p), timings[index])
	}

}

func calculateStats(timings []time.Duration, results chan time.Duration) []time.Duration {
	for r := range results {
		timings = append(timings, r)
	}

	sort.Slice(timings, func(i, j int) bool {
		return timings[i] < timings[j]
	})

	return timings
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

func sliceSum(slice []time.Duration) time.Duration {
	var total time.Duration
	for _, t := range slice {
		total += t
	}
	return total
}
