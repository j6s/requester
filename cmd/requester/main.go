package main;

import (
	"log"
	"strings"
	"flag"
	"sync"
	"fmt"
	"os"
)

var (
	workers = flag.Int("workers", 8, "The number of workers to use")
	listSuccess = flag.Bool("list-success", false, "If specified successful (200) requests will also be printed in the summary")
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, " Requester \n")
		fmt.Fprintf(os.Stderr, "===========\n")
		fmt.Fprintf(os.Stderr, "Simple website crawler that tries to follow every link that points\n")
		fmt.Fprintf(os.Stderr, "back to the same domain.\n\n")

		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "$ ./requester {OPTIONAL ARGUMENTS} http://first-url.com http://second-url.com \n\n")
		fmt.Fprintf(os.Stderr, "Optional arguments:\n")
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()
	scanner := Scanner{}

	// Print help if no URLs given
	if (len(flag.Args()) == 0) {
		flag.Usage()
		return
	}

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
		if status == "200 OK" && !*listSuccess {
			continue
		}
		for _, request := range(requests) {
			log.Printf("\t -> %s", request.url.String())
		}
	}
}
