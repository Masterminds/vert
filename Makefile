DESTDIR=/usr/local/bin

all: bootstrap test install

setup:
	go get -u github.com/golang/dep/cmd/dep

bootstrap:
	dep ensure
	dep status

vert:
	go build -o vert vert.go

test:
	go test .

install: vert
	install -d ${DESTDIR}
	install -m 755 ./vert ${DESTDIR}/vert

.PHONY: bootstrap test install all setup
