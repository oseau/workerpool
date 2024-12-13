include .env

MAKEFLAGS += --always-make
WORKSPACE = /usr/src/workerpool

.DEFAULT_GOAL := help

test: ## Run tests in Docker
	docker run --rm -v $(PWD):$(WORKSPACE) -w $(WORKSPACE) ${GOLANG_IMAGE} go test -count=1 -race -v .

coverage: ## Generate and open test coverage report
	@rm -f .coverage.html
	docker run --rm -v $(PWD):$(WORKSPACE) -w $(WORKSPACE) ${GOLANG_IMAGE} sh -c \
		"go test -coverprofile=.coverage.out . && \
		go tool cover -html=.coverage.out -o .coverage.html && \
		rm .coverage.out"
	@open .coverage.html

bench: ## Run benchmarks in Docker
	docker run --rm -v $(PWD):$(WORKSPACE) -w $(WORKSPACE) ${GOLANG_IMAGE} go test -bench=. -benchmem .

shell: ## Start a shell in the Docker container
	docker run -it --rm -v $(PWD):$(WORKSPACE) -w $(WORKSPACE) ${GOLANG_IMAGE} /bin/sh

example: ## Run the basic example
	docker run --rm -v $(PWD):$(WORKSPACE) -w $(WORKSPACE)/examples/basic ${GOLANG_IMAGE} go run .

help:
	@grep -Eh '^[a-zA-Z_-]+:.*?##? .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?##? "}; {printf "\033[36m%-10s\033[0m %s\n", $$1, $$2}'
