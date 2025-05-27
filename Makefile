.PHONY: build pind

# env vars for all targets
# -static prevents: glibc version not found
.EXPORT_ALL_VARIABLES:
CGO_ENABLED=1
CGO_LDFLAGS="-static"

build:
	go build -v -o ./bin/inact github.com/pnsafonov/inact

build_cross: build_linux build_freebsd

build_linux:
	GOOS=linux \
	go build -v -o ./bin/inact_linux github.com/pnsafonov/inact

build_freebsd:
	GOOS=freebsd \
	go build -v -o ./bin/inact_freebsd github.com/pnsafonov/inact

test:
	go test -count=1 ./...
