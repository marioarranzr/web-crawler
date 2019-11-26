.PHONY: test

test:
	go test -count=1 -cover -v ./...

build:
	go build .

run:
	go run main.go -concurrency $(concurrency) $(url)