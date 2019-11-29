.PHONY: all clean deps build

golint:
	@golangci-lint run ./...
