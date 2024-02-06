package main

import (
	"github.com/gookit/color"
	"golang.org/x/net/html"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type Request struct {
	url        url.URL
	executed   bool
	response   http.Response
	references []Request
}

func (page *Request) Parse(incomingUrl string) {
	uri, err := url.Parse(incomingUrl)
	if err != nil {
		log.Fatal(err)
	}

	// URL Fragments (Everything after the hash) is ignored
	uri.Fragment = ""
	page.url = *uri
}

func (parent Request) Child(childUrl string) *Request {
	if strings.Index(childUrl, "http") != 0 {
		return nil
	}

	if childUrl[0] == '\'' || childUrl[0] == '"' {
		childUrl = childUrl[1:]
	}
	lastIndex := len(childUrl) - 1
	if childUrl[lastIndex] == '\'' || childUrl[lastIndex] == '"' {
		childUrl = childUrl[:lastIndex]
	}

	child := &Request{}
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
	if len(response.Header["Content-Type"]) > 0 && strings.Index(response.Header["Content-Type"][0], "text/html") != -1 {
		for _, href := range extractTagAttribute(response.Body, "a", "href") {
			child := request.Child(href)
			if child != nil {
				request.references = append(request.references, *child)
			}
		}
	}
}

func (request *Request) Matches(comparisonUrl url.URL) bool {
	return request.url.Host == comparisonUrl.Host && strings.Trim(request.url.Path, "/") == strings.Trim(comparisonUrl.Path, "/")
}

func extractTagAttribute(reader io.Reader, tagName string, attributeName string) []string {
	attributes := make([]string, 0)
	tokenizer := html.NewTokenizer(reader)

	for {
		token := tokenizer.Next()
		switch {
		case token == html.ErrorToken:
			// End of document
			return attributes
		case token == html.StartTagToken:
			tag := tokenizer.Token()
			if tag.Data == tagName {
				for _, attr := range tag.Attr {
					if attr.Key == attributeName {
						attributes = append(attributes, strings.TrimSpace(attr.Val))
					}
				}
			}
		}
	}

	return attributes
}
