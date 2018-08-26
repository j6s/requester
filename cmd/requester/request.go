package main;

import (
	"net/http"
	"net/url"
	"log"
	"github.com/gookit/color"
	"strings"
	"io"
	"golang.org/x/net/html"
)

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
