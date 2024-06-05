GO_VERSION := go1.22.4
GC := ~/go/bin/$(GO_VERSION)

all: bin/client bin/server

bin/client:
	mkdir -p bin
	$(GC) build -o ./bin/client ./client

bin/server:
	mkdir -p bin
	$(GC) build -o ./bin/server ./server 

tidy:
	$(GC) mod tidy

install_go:
	go install golang.org/dl/$(GO_VERSION)@latest
	$(GC) download
