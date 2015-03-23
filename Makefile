build:
	go build

release:
	gox -osarch="darwin/amd64 linux/amd64 linux/arm" -output="./bin/musicbot_{{.OS}}_{{.Arch}}"

bootstrap:
	gox -build-toolchain

setup:
	go get github.com/mitchellh/gox
	go get

clean:
	rm -f ./musicbox
	rm -f ./bin/*