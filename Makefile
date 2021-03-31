## Updates the dependencies.
update:
	echo "updating dependencies ..."
	go get -u ./...

## Builds the application for the current OS.
build:
	echo "building bench-routes ..."
	go build -o bench-routes src/main.go

## Runs the UI (assuming all dependencies in dashboard/v1.1 are installed).
view-v1.1:
	cd dashboard/v1.1/ && yarn start

## Installs the UI dependencies, checks for style guide compliance and builds and runs the app.
test-views-v1.1:
	cd dashboard/v1.1/ && yarn install
	cd dashboard/v1.1/ && yarn run lint
	cd dashboard/v1.1/ && yarn run tlint
	cd dashboard/v1.1/ && prettier '**/*.tsx' --list-different
	cd dashboard/v1.1/ && yarn run build
	cd dashboard/v1.1/ && yarn start &

## Removes all residual files.
clean:
	rm -R build/ bench-routes

## Runs Golang unit tests
test: build
	go clean -testcache
	go test -v ./...

## Installs UI dependencies and builds the frontend.
build-frontend:
	cd dashboard/v1.1/ && yarn install --network-timeout 1000000 && yarn build

## Runs Golang unit tests without mentioning the skipped tests.
test-non-verbose: build
	go clean -testcache
	go test ./...

## Runs selenium tests.
test-services: build
	./bench-routes &
	cd tests && yarn install
	yarn global add mocha
	mocha tests/browser.js

## Complete testing include building for all supported OS.
test_complete: build
	./shell/go-build-all.sh
	echo "test success!"

## Executes the application (assuming all dependencies are installed)
run:
	echo "compiling go-code and executing bench-routes"
	echo "using 9990 as default service listener port"
	go run src/*.go 9990

run-collector:
	go run src/collector/main.go

## Fixes webapp and server code style.
fix:
	go fmt ./...
	cd dashboard/v1.1/ && npm run prettier-fix
	cd dashboard/v1.1/ && npm run tlint-fix

## Runs golangci-lint (assuming golangci-lint is installed).
lint:
	@if ! [ -x "$$(command -v golangci-lint)" ]; then \
		echo "golangci-lint is not installed. Please see https://github.com/golangci/golangci-lint#install for installation instructions."; \
		exit 1; \
	fi;
	golangci-lint run

# Help documentation à la https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help:
	@cat Makefile | grep -v '\.PHONY' |  grep -v '\help:' | grep -B1 -E '^[a-zA-Z0-9_.-]+:.*' | sed -e "s/:.*//" | sed -e "s/^## //" |  grep -v '\-\-' | sed '1!G;h;$$!d' | awk 'NR%2{printf "\033[36m%-30s\033[0m",$$0;next;}1' | sort
