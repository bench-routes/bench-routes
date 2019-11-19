update:
	echo "updating dependencies ..."
	go get -u ./...

build:
	echo "building bench-routes ..."
	go build src/main.go src/handlers.go
	mv main bench-routes

clean:
	rm -R build/ bench-routes

test: build
	go clean -testcache
	go test -v ./...

test-services: build
	./bench-routes &
	cd tests && npm install
	npm install -g mocha
	mocha tests/browser.js

test_complete: build
	./shell/go-build-all.sh
	echo "test success!"

run:
	echo "compiling go-code and executing bench-routes"
	echo "using 9090 as default service listener port"
	go run src/*.go 9090

fix:
	go fmt ./...

lint:
	golangci-lint run

