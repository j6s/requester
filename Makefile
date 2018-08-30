build:
	go build -o requester ./cmd/requester/*.go

clean:
	rm -Rfv bin
	mkdir bin

build: clean
	go build -o bin/requester ./cmd/requester/

build-all: clean
	GOOS="linux"   GOARCH="amd64"       go build -o bin/requester__linux-amd64 ./cmd/requester/
	GOOS="linux"   GOARCH="arm" GOARM=6 go build -o bin/requester__linux-armv6 ./cmd/requester/
	GOOS="linux"   GOARCH="arm" GOARM=7 go build -o bin/requester__linux-armv7 ./cmd/requester/
	GOOS="linux"   GOARCH="arm"         go build -o bin/requester__linux-arm   ./cmd/requester/
	GOOS="darwin"  GOARCH="amd64"       go build -o bin/requester__macos-amd64 ./cmd/requester/
	GOOS="windows" GOARCH="amd64" go build -o bin/requester__win-amd64 *.go
