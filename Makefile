GC := go1.22.4

install_go:
	go install golang.org/dl/$(GC)@latest
	$(GC) download
