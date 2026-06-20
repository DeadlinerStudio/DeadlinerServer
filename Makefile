MODULE := github.com/aritxonly/deadlinerserver
SERVICE := deadliner
IDL := idl/deadliner.thrift

.PHONY: fmt test test-rpc run generate

fmt:
	gofmt -w ./main.go ./handler.go ./internal

test:
	GOCACHE=/private/tmp/deadlinerserver-gocache go test -vet=off ./internal/app ./internal/config ./internal/domain/... ./internal/utils

test-rpc:
	GOCACHE=/private/tmp/deadlinerserver-gocache go test -vet=off ./...

run:
	GOCACHE=/private/tmp/deadlinerserver-gocache go run .

generate:
	kitex -module $(MODULE) -service $(SERVICE) $(IDL)
