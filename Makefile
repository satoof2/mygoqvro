NAME=phybbit
LDFLAGS := -X main.revision=${REVISION}

GOCMD=go
GOBUILD=$(GOCMD) build
GOFMT=goimports

build-linux:
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o bin/linux/$(BINARY) -ldflags "$(LDFLAGS)" src/main.go
