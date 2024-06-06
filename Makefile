GO_VERSION := go1.22.4
ifeq ($(OS), Windows_NT)
	GC := go
else
	GC := ~/go/bin/$(GO_VERSION)
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

tidy:
	$(GC) mod tidy

install_go:
	go install golang.org/dl/$(GO_VERSION)@latest
	$(GC) download

bin:
	mkdir bin

FORCE: ;