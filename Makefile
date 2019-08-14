BIN_DIR:=.bin

.PHONY: help
help: ## Show this help
	@echo "Execute one of this targets: "
	@echo
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/:.*##/:##/' | column -t -s '##'

.PHONY: go-tools-install
go-tools-install: .mkdir-bin ## Install Go tools
	@GOBIN=$(realpath $(BIN_DIR)) go install github.com/golangci/golangci-lint/cmd/golangci-lint
	@GOBIN=$(realpath $(BIN_DIR)) go install github.com/mitchellh/gox

.PHONY: build-bins
build-bins: ## Build the binaries
	@mkdir -p build
	@go build -o build/dbcleaner ./dbcleaner

.PHONY: lint
lint: ## Lint the sources of all the packages contained in this repo
	@PATH=$(realpath $(BIN_DIR)):$$PATH golangci-lint run --enable-all --exclude-use-default=false

.PHONY: test
test: ## Execute all the tests
	@go test -race $(TARGS) ./...

.PHONY: .mkdir-bin
.mkdir-bin:
	@mkdir -p .bin
