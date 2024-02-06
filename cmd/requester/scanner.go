package main

type Scanner struct {
	pool         []*Request
	queue        []*Request
	allowedHosts []string
}

func (scanner *Scanner) isHostAllowed(host string) bool {
	for _, allowed := range scanner.allowedHosts {
		if allowed == host {
			return true
		}
	}
	return false
}

func (scanner *Scanner) isAlreadyInQueue(url string) bool {
	for _, request := range scanner.pool {
		if request.url.String() == url {
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
	for _, child := range request.references {
		scanner.PushToQueue(child)
	}

	return true
}

func (scanner *Scanner) RequestsByStatus() map[string][]Request {
	requests := make(map[string][]Request, 0)
	for _, req := range scanner.pool {
		_, exists := requests[req.response.Status]
		if !exists {
			requests[req.response.Status] = make([]Request, 0)
		}

		requests[req.response.Status] = append(requests[req.response.Status], *req)
	}

	return requests
}

func (scanner *Scanner) Requests() []*Request {
	return scanner.pool
}
