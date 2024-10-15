.DEFAULT_GOAL := goapp client

.PHONY: all
all: clean goapp client

.PHONY: goapp
goapp:
	mkdir -p bin
	go build -o bin/server ./cmd/server

.PHONY: client
client:
	mkdir -p bin
	go build -o bin/client ./cmd/client

.PHONY: clean
clean:
	go clean
	rm -f bin/*