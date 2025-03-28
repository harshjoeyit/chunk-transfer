package main

import (
	"io"
	"net/http"
	"testing"
)

const (
	baseURL     = "http://localhost:8080"
	imagePaths  = "/images/timg1.png,/images/timg2.png,/images/timg3.png,/images/timg4.png,/images/timg5.png"
	concurrency = 10 // Number of concurrent requests for the concurrent benchmark
)

// fetchThumbnails fetches thumbnails from the given URL.
func fetchThumbnails(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Read the entire response body (important for chunked encoding)
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return err
	}

	return nil
}

// Benchmark - concurrent non-blocking fetching of thumbnails
// call API - http://localhost:8080/thumbnail-batch-concurrent?paths=/images/timg1.png,/images/timg2.png,/images/timg3.png,/images/timg4.png,/images/timg5.png
func BenchmarkConcurrentWithGoroutines(b *testing.B) {
	apiURL := baseURL + "/thumbnail-batch-concurrent?paths=" + imagePaths

	b.SetParallelism(concurrency)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if err := fetchThumbnails(apiURL); err != nil {
				b.Fatalf("Error fetching thumbnails: %v", err)
			}
		}
	})
}

func BenchmarkConcurrent(b *testing.B) {
	apiURL := baseURL + "/thumbnail-batch-concurrent?paths=" + imagePaths

	for n := 0; n < b.N; n++ {
		// Use a WaitGroup to wait for the goroutine to complete
		// Launch a goroutine to fetch the thumbnails
		if err := fetchThumbnails(apiURL); err != nil {
			b.Fatalf("Error fetching thumbnails: %v", err)
		}
	}
}

// Benchmark - Blocking (blocking) fetching of thumbnails
// call API - http://localhost:8080/thumbnail-batch-blocking?paths=/images/timg1.png,/images/timg2.png,/images/timg3.png,/images/timg4.png,/images/timg5.png
func BenchmarkBlockingWithGoroutines(b *testing.B) {
	apiURL := baseURL + "/thumbnail-batch-blocking?paths=" + imagePaths

	b.SetParallelism(concurrency)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if err := fetchThumbnails(apiURL); err != nil {
				b.Fatalf("Error fetching thumbnails: %v", err)
			}
		}
	})
}

func BenchmarkBlocking(b *testing.B) {
	apiURL := baseURL + "/thumbnail-batch-blocking?paths=" + imagePaths

	for n := 0; n < b.N; n++ {
		// Use a WaitGroup to wait for the goroutine to complete
		// Launch a goroutine to fetch the thumbnails
		if err := fetchThumbnails(apiURL); err != nil {
			b.Fatalf("Error fetching thumbnails: %v", err)
		}
	}
}
