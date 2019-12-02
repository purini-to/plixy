.PHONY: all clean deps build

golint:
	@golangci-lint run ./...

test:
	@go test -v -cover ./...

build:
	@go build -v .
