package main;

import (
	"log"
	"strings"
	"flag"
)

var (
	workers = flag.Int("workers", 8, "The number of workers to use")
)


func main() {
	flag.Parse()
	scanner := Scanner{}

	// Initialize parser with command line arguments
	for _, base := range(flag.Args()) {
		if strings.Index(base, "http") == -1 {
			base = "http://" + base
		}
		scanner.Initialize(base)
	}

	// Setup thread pool
	var wg sync.WaitGroup
	for w := 1; w <= *workers; w++ {
		wg.Add(1)
		go func() {
			for {
				if !scanner.Work() {
					break
				}
			}
			wg.Done()
		}()
	}

	// Wait for All workers to terminate
	wg.Wait()

	// Print stats
	for status, requests := range(scanner.RequestsByStatus()) {
		log.Printf("## %v: %d Requests", status, len(requests))
		if status == "200 OK" {
			continue
		}
		for _, request := range(requests) {
			log.Printf("\t -> %s", request.url.String())
		}
	}
}
