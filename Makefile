Version := $(shell git describe --tags --dirty 2> /dev/null)
GitCommit := $(shell git rev-parse HEAD)
LDFLAGS := "-s -w -X github.com/tschaefer/rpinfo/version.Version=$(Version) -X github.com/tschaefer/rpinfo/version.GitCommit=$(GitCommit)"

.PHONY: all
all: fmt lint test dist

.PHONY: fmt
fmt:
	test -z $(shell gofmt -l .) || (echo "[WARN] Fix format issues" && exit 1)

.PHONY: lint
lint:
	test -z $(shell golangci-lint run >/dev/null || echo 1) || (echo "[WARN] Fix lint issues" && exit 1)

.PHONY: dist
dist:
	mkdir -p bin
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o bin/rpinfo-arm64 -ldflags $(LDFLAGS) .
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=5 go build -o bin/rpinfo-armv5 -ldflags $(LDFLAGS) .
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 go build -o bin/rpinfo-armv6 -ldflags $(LDFLAGS) .
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -o bin/rpinfo-armv7 -ldflags $(LDFLAGS) .

.PHONY: checksum
checksum:
	cd bin && \
	for f in rpinfo-arm64 rpinfo-armv5 rpinfo-armv6 rpinfo-armv7; do \
		sha256sum $$f > $$f.sha256; \
	done && \
	cd ..

.PHONY: test
test:
	test -z $(shell go test ./... 2>&1 >/dev/null || echo 1) || (echo "[WARN] Fix test issues" && exit 1)

.PHONY: clean
clean:
	rm -rf bin
