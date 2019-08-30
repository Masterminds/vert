all: test install

build:
	go build -o vert vert.go

test:
	go test .

install: build
	install -d ${DESTDIR}/usr/local/bin/
	install -m 755 ./vert ${DESTDIR}/usr/local/bin/vert

.PHONY: test build install all
