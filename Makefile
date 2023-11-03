GO ?= go
GOFMT ?= gofmt "-s"
GO_VERSION=$(shell $(GO) version | cut -c 14- | cut -d' ' -f1 | cut -d'.' -f2)
PACKAGES ?= $(shell $(GO) list ./...)
GO_FILES := $(shell find . -name "*.go" -not -path "./vendor/*" -not -path ".git/*")
TEST_TAGS ?= ""
GIT_COMMIT_SHA := $(shell git rev-parse HEAD | cut -c 1-8)

BINARY_NAME=toolbox

dep:
	@go mod tidy && go fmt

build:
	@GOARCH=arm64 GOOS=darwin go build -o bin/${BINARY_NAME}_darwin_arm64 main.go

build-race:
	@GOARCH=arm64 GOOS=darwin go build -race -o bin/${BINARY_NAME}_darwin_arm64 main.go


run: build
	@echo $(MAKECMDGOALS)
	./bin/${BINARY_NAME}_darwin_arm64

dev: 
	@CompileDaemon --command="./bin/daemon-build" -build="go build -o bin/daemon-build" -color=true -graceful-kill=true

clean:
	go clean
	rm bin/${BINARY_NAME}_linux_arm64

lint:
	@hash golint > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		$(GO) get -u golang.org/x/lint/golint; \
	fi
	for PKG in $(PACKAGES); do golint -set_exit_status $$PKG || exit 1; done;

install-tools:
	if [ $(GO_VERSION) -gt 15 ]; then \
		$(GO) install golang.org/x/lint/golint@latest; \
		$(GO) install github.com/client9/misspell/cmd/misspell@latest; \
	elif [ $(GO_VERSION) -lt 16 ]; then \
		$(GO) install golang.org/x/lint/golint; \
		$(GO) install github.com/client9/misspell/cmd/misspell; \
	fi
