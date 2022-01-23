GOLANG_CI_VERSION ?= v1.43.0

GIT_HASH := $(shell git rev-parse --short HEAD)
LDFLAGS := -ldflags="-s -w -X main.version=$(GIT_HASH)"

.PHONY: clean
clean:
	@rm -rf dist

.PHONY: build
build:
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -trimpath -o ./dist/ ./cmd/...

.PHONY: lint
lint:
	@docker run --rm -v $(shell pwd):/app -w /app -it golangci/golangci-lint:$(GOLANG_CI_VERSION) golangci-lint run

.PHONY: test
test:
	@go test -v -cover ./...