BINARY := ghostio

.PHONY: build test clean

build:
	go build -o $(BINARY) ./cmd/ghostio

test:
	go test ./...

clean:
	rm -f $(BINARY)
