.PHONY: test

build:
	go build -o dist/nv main.go

run:
	go run main.go

test: build
	go test -v ./...

uninstall:
	rm /usr/local/bin/nv

release:
	sh ./bin/release.sh
