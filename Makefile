NAME=phybbit
REVISION:=$(shell git rev-parse --short HEAD)
LDFLAGS := -X main.revision=${REVISION}

GOCMD=go
GOBUILD=$(GOCMD) build
GOFMT=goimports

build-linux:
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(NAME) -ldflags "-X main.mode=encoder" cmd/main.go
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o decoder -ldflags "-X main.mode=decoder" cmd/main.go

build:
	$(GOBUILD) -o $(NAME) -ldflags "-X main.mode=encoder" cmd/main.go
	$(GOBUILD) -o decoder -ldflags "-X main.mode=decoder" cmd/main.go
