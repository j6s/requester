package main;

import (
	"net/http"
	"log"
	"io"
	"golang.org/x/net/html"
	"strings"
	"net/url"
	"github.com/gookit/color"
	"flag"
)

type Scanner struct {
	pool []*Request;
	queue []*Request;
	allowedHosts []string;
}

func (scanner *Scanner) isHostAllowed(host string) bool {
	for _, allowed := range(scanner.allowedHosts) {
		if allowed == host {
			return true
		}
	}
	return false
}

func (scanner *Scanner) isAlreadyInQueue(url string) bool {
	for _, request := range(scanner.pool) {
		if (request.url.String() == url) {
			return true
		}
	}
	return false
}

func (scanner *Scanner) isRequestAllowed(request Request) bool {
	return (request.url.Scheme == "http" || request.url.Scheme == "https") &&
	 scanner.isHostAllowed(request.url.Host) &&
	!scanner.isAlreadyInQueue(request.url.String())
}

// Checks if the given request is elegible for the queue
// (if the host is allowed and it has not already been processed)
// and pushes it to the queue if it is.
func (scanner *Scanner) PushToQueue(request Request) {
	if scanner.isRequestAllowed(request) {
		scanner.pool = append(scanner.pool, &request)
		scanner.queue = append(scanner.queue, &request)
	}
}

func (scanner *Scanner) Initialize(firstUrl string) {
	request := Request{}
	request.Parse(firstUrl)
	scanner.allowedHosts = append(scanner.allowedHosts, request.url.Host)
	scanner.PushToQueue(request)
}

func (scanner *Scanner) Pop() *Request {
	request := scanner.queue[0]
	scanner.queue = scanner.queue[1:]
	return request
}

func (scanner *Scanner) HasQueuedItems() bool {
	return len(scanner.queue) > 0
}

func (scanner *Scanner) Work() bool {
	if !scanner.HasQueuedItems() {
		return false
	}

	request := scanner.Pop()
	request.Execute()
	for _, child := range(request.references) {
		scanner.PushToQueue(child)
	}
	return true
}

func (scanner *Scanner) RequestsByStatus() map[string][]Request {
	requests := make(map[string][]Request, 0)
	for _, req := range(scanner.pool) {
		_, exists := requests[req.response.Status]
		if !exists {
			requests[req.response.Status] = make([]Request, 0)
		}

		requests[req.response.Status] = append(requests[req.response.Status], *req)
	}

	return requests;
}

type Request struct {
	url url.URL;
	executed bool;
	response http.Response;
	references []Request;
}

func (page *Request) Parse(incomingUrl string) {
	uri, err := url.Parse(incomingUrl)
	if err != nil {
		log.Fatal(err)
	}
	page.url = *uri
}

func (parent Request) Child(childUrl string) Request {
	child := Request{}
	child.Parse(childUrl)
	if child.url.Host == "" {
		child.url.Host = parent.url.Host
	}
	if child.url.Scheme == "" {
		child.url.Scheme = parent.url.Scheme
	}
	return child
}

func (request *Request) Execute() {
	response, err := http.Get(request.url.String())
	if err != nil {
		request.response.StatusCode = 999
		request.response.Status = "999 Something happened"
	} else {
		request.response = *response
	}

	c := color.New(color.FgGreen)
	if request.response.StatusCode > 200 && request.response.StatusCode < 400 {
		c = color.New(color.FgYellow)
	} else if request.response.StatusCode > 400 {
		c = color.New(color.FgWhite, color.BgRed)
	}
	log.Printf(c.Render("[%v] GET %v"), response.Status, request.url.String())


	request.executed = true
	if strings.Index(response.Header["Content-Type"][0], "text/html") != -1 {
		for _, href := range(extractTagAttribute(response.Body, "a", "href")) {
			request.references = append(request.references, request.Child(href))
		}
	}
}


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

func extractTagAttribute(reader io.Reader, tagName string, attributeName string) []string {
	attributes := make([]string, 0)
	tokenizer := html.NewTokenizer(reader)

	for {
		token := tokenizer.Next()
		switch {
		case token == html.ErrorToken:
			// End of document
			return attributes;
		case token == html.StartTagToken:
			tag := tokenizer.Token()
			if tag.Data == tagName {
				for _, attr := range tag.Attr {
					if attr.Key == attributeName {
						attributes = append(attributes, attr.Val)
					}
				}
			}
		}
	}

	return attributes;
}
