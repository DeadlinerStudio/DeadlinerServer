MODULE := github.com/aritxonly/deadlinerserver
SERVICE := deadliner
IDL := idl/deadliner.thrift
TMP_ROOT ?= $(if $(TMPDIR),$(TMPDIR),/tmp)
GO_CACHE_DIR ?= $(TMP_ROOT)/deadlinerserver-gocache

.PHONY: fmt test test-rpc run package generate

fmt:
	gofmt -w ./cmd ./internal

test:
	GOCACHE=$(GO_CACHE_DIR) go test -vet=off ./internal/app ./internal/config ./internal/domain/... ./internal/utils

test-rpc:
	GOCACHE=$(GO_CACHE_DIR) go test -vet=off ./...

run:
	GOCACHE=$(GO_CACHE_DIR) go run ./cmd/deadlinerserver

package:
	./script/build.sh

generate:
	kitex -module $(MODULE) -service $(SERVICE) $(IDL)
