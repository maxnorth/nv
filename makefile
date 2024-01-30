
build:
	go build -o dist/nv main.go

run:
	go run main.go

uninstall:
	rm /usr/local/bin/nv

release:
	sh ./bin/release.sh
