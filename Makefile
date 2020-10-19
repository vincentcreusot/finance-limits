pkgs   = $(shell go list ./...)

test: build ## running test after build
	@echo ">> running tests"
	@go test -v -short $(pkgs)

coverage: build ## gives test coverage
	@go test -short -cover $(pkgs)

format: ## Format code
	@echo ">> formatting code"
	@go fmt $(pkgs)

vet: ## vet code
	@echo ">> vetting code"
	@go vet $(pkgs)

build: ## build code with promu
	@echo ">> building binaries"
	@go build

lint: golint ## lint code
	@echo ">> linting code"
	@! golint $(pkgs) | grep '^'

golint: ## gets golint for building
	@go get -u golang.org/x/lint/golint
