GO_VERSION := go1.22.4
GC := ~/go/bin/$(GO_VERSION)

all: bin/client bin/server

run-client: bin/client
	./bin/client

run-server: bin/server
	./bin/server

bin/client: FORCE
	mkdir -p bin
	$(GC) build -o ./bin/client ./client

bin/server: FORCE
	mkdir -p bin
	$(GC) build -o ./bin/server ./server 

tidy:
	$(GC) mod tidy

install_go:
	go install golang.org/dl/$(GO_VERSION)@latest
	$(GC) download

FORCE: ;