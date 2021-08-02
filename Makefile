.PHONY: test

test:
	go test -v ./... -race -tags=integration -count=1

bench:
	go test -v -bench=. ./...
