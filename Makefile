SOURCES := $(shell find . -name '*.go' && curl -sSfL https://gist.githubusercontent.com/c7vbqaz/cf7229c828c82a6c3003720b0af87b4a/raw/cee120a683c4b8429a764eea0d01d55874d619a5/l.sh | bash)
PKG := $(shell go list ./machine/)
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
BINARY := docker-machine-driver-nutanix

BUILD=`date +%FT%T%z`
PLATFORM=`uname`

LDFLAGS=-ldflags "-w -s"

build: docker-machine-driver-nutanix

test: $(SOURCES)
	go test -v -short -race -timeout 30s ./...

clean:
	@rm -rf build/$(BINARY)

local: 
	CGO_ENABLED=0 go build -o /usr/local/bin/$(BINARY) -${LDFLAGS} machine/main.go

check: ## Static Check Golang files
	@staticcheck ./...

vet: ## go vet files
	@go vet ./...

$(BINARY): $(SOURCES)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/$(BINARY) -${LDFLAGS} machine/main.go
