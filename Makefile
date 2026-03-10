.PHONY: build install test clean

PREFIX ?= /usr/local/bin

build:
	go build -o towr ./cmd/towr/

install: build
	install -m 755 towr $(PREFIX)/towr

test:
	go test ./... -count=1

clean:
	rm -f towr
