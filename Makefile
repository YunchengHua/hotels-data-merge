.PHONY: help
help: ## Show list of make targets and their description
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| sort \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: test
test: ## Run tests
	go test -v -cover ./...

.PHONY: lint
lint: ## Run linter
	golangci-lint run

.PHONY: build
build: ## Build Docker image
	docker build -t hotels-data-merge .

.PHONY: run
run: ## Run Docker image
	docker run -p 8080:8080 hotels-data-merge

.PHONY: log
log: ## Show application logs from Docker container
	@CONTAINER_ID=$$(docker ps | grep hotels-data-merge | awk '{print $$1}'); \
	if [ -z "$$CONTAINER_ID" ]; then \
		echo "No running hotels-data-merge container found."; \
	else \
		docker exec $$CONTAINER_ID cat /log/application.log; \
	fi
