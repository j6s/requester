package main;

import (
	"log"
	"strings"
	"flag"
)

func main() {
	flag.Parse()
	scanner := Scanner{}

	log.Printf("%#v", flag.Args())
	for _, base := range(flag.Args()) {
		if strings.Index(base, "http") == -1 {
			base = "http://" + base
		}
		scanner.Initialize(base)
	}

	log.Printf("%#v", scanner)
	for {
		if !scanner.Work() {
			break
		}
	}

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
