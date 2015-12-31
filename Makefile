all: bootstrap test install

bootstrap:
	glide install

build:
	go build -o vert vert.go

test:
	go test .

install: build
	install -d ${DESTDIR}/usr/local/bin/
	install -m 755 ./vert ${DESTDIR}/usr/local/bin/vert

.PHONY: bootstrap test build install all
