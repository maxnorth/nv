
build:
	go build -o bin/nv main.go

run:
	go run main.go

release:
	sh ./bin/release.sh
