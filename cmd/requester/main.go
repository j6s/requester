package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"sync"
)

var (
	workers       = flag.Int("workers", 8, "The number of workers to use")
	listSuccess   = flag.Bool("list-success", false, "If specified successful (200) requests will also be printed in the summary")
	diffToSitemap = flag.String("diff-to-sitemap", "", "If specified, a list of discovered URLs not in the specified sitemap will be shown.")
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
	if len(flag.Args()) == 0 {
		flag.Usage()
		return
	}

	// Initialize parser with command line arguments
	for _, base := range flag.Args() {
		if strings.Index(base, "http") == -1 {
			base = "http://" + base
		}
		scanner.Initialize(base)
	}

	// Initialize sitemap
	var urlsInSitemap []url.URL
	var err error
	if *diffToSitemap != "" {
		urlsInSitemap, err = GetUrlsInSitemap(*diffToSitemap)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Setup thread pool
	var wg sync.WaitGroup
	for w := 1; w <= *workers; w++ {

		// Execute one unit of synchronous work per
		// worker being used. This fills the queue enough
		// for workers not to starve each other.
		if !scanner.Work() {
			break
		}

		// Start worker.
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
	for status, requests := range scanner.RequestsByStatus() {
		log.Printf("## %v: %d Requests", status, len(requests))
		if status == "200 OK" && !*listSuccess {
			continue
		}
		for _, request := range requests {
			log.Printf("\t -> %s", request.url.String())
		}
	}

	if len(urlsInSitemap) > 0 {
		for _, request := range scanner.Requests() {
			isInSitemap := false
			for _, urlInSitemap := range urlsInSitemap {
				if request.Matches(urlInSitemap) {
					isInSitemap = true
					break
				}
			}

			if !isInSitemap {
				fmt.Printf("\n- %s", request.url.String())
			}
		}
	}
}
