all: bootstrap test install

setup:
	go get -u github.com/golang/dep/cmd/dep

bootstrap:
	dep ensure
	dep status

build:
	go build -o vert vert.go

test:
	go test .

install: build
	install -d ${DESTDIR}/usr/local/bin/
	install -m 755 ./vert ${DESTDIR}/usr/local/bin/vert

.PHONY: bootstrap test build install all setup
