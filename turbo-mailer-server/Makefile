GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=target/turbo-mailer

.PHONY: all
all: build

.PHONY: test
test:
	$(GOTEST) -v ./...

.PHONY: clean
clean:
	$(GOCLEAN)
	rm -rf $(BINARY_NAME)

.PHONY: deps
deps:
	$(GOGET) -v -t -d ./...
	$(GOMOD) tidy

.PHONY: generate
generate:
	$(GOCMD) generate ./...

.PHONY: build
build: generate
	$(GOBUILD) -o $(BINARY_NAME) -v -ldflags "-s -w"

.PHONY: build-linux
build-linux: generate
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME)_linux -v -ldflags "-s -w"

.PHONY: build-windows
build-windows: generate
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME).exe -v -ldflags "-s -w"

.PHONY: build-darwin
build-darwin: generate
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME)_darwin -v -ldflags "-s -w"

.PHONY: dev
dev:
	@echo "Starting development server..."
	@which air > /dev/null || go install github.com/air-verse/air@latest
	air -c .air.toml
