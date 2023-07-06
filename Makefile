# TODO: add Docker-related targets

build:
	go build -o kagglecarapi

run: build
	./kagglecarapi

test:
	go test -v ./...