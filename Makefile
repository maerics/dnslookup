.PHONY: test fmt vet modtidy build

test: fmt vet
	go test ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

modtidy:
	go mod tidy

build:
	go build -o dnslookup .
