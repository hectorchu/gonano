.ONESHELL:
.PHONY: build

build:
	go build .

test:
	go test -v -cover ./...

lint:
	@command -v golangci-lint > /dev/null || \
		GO111MODULE=off go get -u github.com/golangci/golangci-lint/cmd/golangci-lint

	golangci-lint run

