default: clear build-osx build-linux

test:
	go test -v ./... -count=1

clear:
	echo "done"

build-osx:
	go build -o bin/osx/go-watcher cmd/go-watcher/main.go

build-linux:
	GOOS=linux GOARCH=amd64 go build -o bin/linux/go-watcher cmd/go-watcher/main.go