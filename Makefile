.PHONY: build test test-coverage lint fmt vet clean install all

BINARY    := i18n-fixer
VERSION   := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT    := $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE      := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS   := -ldflags "-s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)"

build:
	go build $(LDFLAGS) -o bin/$(BINARY) ./cmd/i18n-fixer/

test:
	go test -race ./...

test-coverage:
	go test -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

lint:
	golangci-lint run

fmt:
	gofmt -s -w .

vet:
	go vet ./...

clean:
	rm -rf bin/ coverage.out coverage.html dist/

install:
	go install $(LDFLAGS) ./cmd/i18n-fixer/

all: fmt vet test build
