GO_VERSION := go1.22.4
GO_PATH := ~/go/bin/$(GO_VERSION)

ifeq ($(shell test -s $(GO_PATH) && echo -n yes), yes)
	GC := $(GO_PATH)
else
	GC := go
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

tidy:
	$(GC) mod tidy

install_go:
	go install golang.org/dl/$(GO_VERSION)@latest
	$(GC) download

bin:
	mkdir bin

bin/linux-arm:
	mkdir bin/linux-arm

FORCE: ;
