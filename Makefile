.DEFAULT_GOAL := all

.PHONY: all
all: clean goapp client

.PHONY: create_bin_folder
create_bin_folder:
	mkdir -p bin

.PHONY: goapp
goapp: create_bin_folder
	go build -o bin/server ./cmd/server

.PHONY: client
client:
	go build -o bin/client ./cmd/client

.PHONY: clean
clean:
	go clean
	rm -f bin/*