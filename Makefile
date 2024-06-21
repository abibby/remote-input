GO_VERSION := go1.22.4
GO_PATH := ~/go/bin/$(GO_VERSION)

ifeq ("$(wildcard $(GO_PATH))", "")
	GC := go
else
	GC := $(GO_PATH)
endif

all: bin/client bin/server

run-client: bin/client
	./bin/client

run-server: bin/server
	./bin/server

bin/client: FORCE bin
	$(GC) build -o ./bin/client ./client

bin/server: FORCE bin
	$(GC) build -o ./bin/server ./server

bin/linux-arm64/server: FORCE bin/linux-arm
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 $(GC) build -o ./bin/linux-arm64/server ./server 

bin/windows-amd64/client.exe: FORCE bin/windows-amd64
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GC) build -o ./bin/windows-amd64/client.exe ./client 

tidy:
	$(GC) mod tidy

install_go:
	go install golang.org/dl/$(GO_VERSION)@latest
	$(GC) download

bin:
	mkdir bin

bin/linux-arm:
	mkdir bin/linux-arm
bin/windows-amd64:
	mkdir bin/windows-amd64

FORCE: ;
