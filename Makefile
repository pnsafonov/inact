.PHONY: build pind

build:
	go build -v -o ./bin/inact github.com/pnsafonov/inact

build_cross: build_linux build_freebsd

build_linux:
	env GOOS=linux \
	CGO_ENABLED=1 \
	go build -v -o ./bin/inact_linux github.com/pnsafonov/inact

build_freebsd:
	env GOOS=freebsd \
	CGO_ENABLED=1 \
	go build -v -o ./bin/inact_freebsd github.com/pnsafonov/inact

test:
	go test -count=1 ./...
