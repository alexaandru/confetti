all: fmt lint actionlint vulncheck deadcode test

test:
	@go test -vet all -coverprofile=unit.cov -covermode=atomic -race -count=5 $(OPTS) ./...
	@go tool cover -func=unit.cov|tail -n1
	@go tool -modfile=tools/go.mod stampli -quiet -coverage=$$(go tool cover -func=unit.cov|tail -n1|tr -s "\t"|cut -f3|tr -d "%")

lint:
	@go run golang.org/x/tools/gopls/internal/analysis/modernize/cmd/modernize@latest -test ./...
	@go tool -modfile=tools/go.mod golangci-lint config verify
	@go tool -modfile=tools/go.mod golangci-lint run

lint-%:
	@go tool -modfile=tools/go.mod golangci-lint --enable-only="$(patsubst lint-%,%,$@)" run

actionlint:
	@go tool -modfile=tools/go.mod actionlint $(OPTS)

vulncheck:
	@go tool -modfile=tools/go.mod govulncheck ./...

deadcode:
	@go tool -modfile=tools/go.mod deadcode -test ./...

fmt:
	@go fmt ./...
	@go tool -modfile=tools/go.mod goimports -l -w .
	@go run mvdan.cc/gofumpt@v0.8.0 -l -w -extra .

doc:
	@go tool -modfile=tools/go.mod godoc -http=:6060 &
	@xdg-open http://localhost:6060/pkg/github.com/alexaandru/confetti/

coverage_map: test
	@go tool -modfile=tools/go.mod go-cover-treemap -coverprofile unit.cov > unit.svg

clean:
	@rm *.cov
	@killall -q godoc || true
