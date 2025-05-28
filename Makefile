all: fmt lint test

fmt:
	@find -name "*.go"|GOFLAGS=-tags=testaws xargs go tool -modfile=tools/go.mod gofumpt -extra -w
	@find -name "*.go"|GOFLAGS=-tags=testaws xargs go tool -modfile=tools/go.mod goimports -w

lint:
	@printf "Linter: "
	@go tool -modfile=tools/go.mod golangci-lint config verify
	@go tool -modfile=tools/go.mod golangci-lint run
	@#go tool -modfile tools/go.mod modernize -test ./...

test:
	@go test -vet all -coverprofile=unit.cov -covermode=atomic $(OPTS) ./...
	@go tool cover -func=unit.cov|tail -n1

testaws:
	@AWS_PROFILE=alex GOFLAGS=-tags=testaws make -s test

doc:
	@go tool -modfile=tools/go.mod godoc -http=:8080 -index

coverage_map:
	@make -s test OPTS=-tags=testaws
	@go tool -modfile=tools/go.mod go-cover-treemap -coverprofile unit.cov > unit.svg

.PHONY: test
