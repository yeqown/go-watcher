versionCode=2.0.1

default: clear compile-osx compile-linux archived

test:
	go test -v ./... -count=1

clear:
	rm -fr package

compile-osx: version
	go build -o package/osx/go-watcher cmd/go-watcher/main.go
	

compile-linux: version
	GOOS=linux GOARCH=amd64 go build -o package/linux/go-watcher cmd/go-watcher/main.go
	

archived:
	- mkdir -p package/archived
	cd package/osx && tar -zcvf ../archived/go-watcher.osx.tar.gz .
	cd package/linux && tar -zcvf ../archived/go-watcher.linux.tar.gz .

version:
	- mkdir -p package/osx
	- mkdir -p package/linux
	echo "${versionCode}" > VERSION
	cp VERSION package/osx
	cp VERSION package/linux