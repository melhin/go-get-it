
.PHONY: help

help: ## This help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

pg_up:
	docker-compose up -d db

pg_down:
	docker-compose stop db

run_test:
	go clean -testcache && go test -run "^Test*" go-get-it/tests -v

test: pg_up run_test pg_down

run: pg_up
	go run main.go