
build:
	go build -o dist/nv main.go

run:
	go run main.go

release:
	sh ./bin/release.sh
