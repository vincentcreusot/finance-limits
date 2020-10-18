pkgs   = $(shell go list ./...)

test: build ## running test after build
	@echo ">> running tests"
	@go test -v -short $(pkgs)

cover: build
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