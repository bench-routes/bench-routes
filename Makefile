update:
	echo "updating dependencies ..."
	go get -u ./...

build:
	echo "building bench-routes ..."
	go build src/main.go src/handlers.go
	mv main bench-routes

view:
	cd dashboard/v1.0/ && sudo npm start

test-views:
	cd dashboard/v1.0/ && npm install
	cd dashboard/v1.0/ && npm run lint
	cd dashboard/v1.0/ && npm run tlint
	cd dashboard/v1.0/ && prettier '**/*.tsx' --list-different
	cd dashboard/v1.0/ && npm run react-test
	cd dashboard/v1.0/ && npm run react-build
	cd dashboard/v1.0/ && npm run build
	cd dashboard/v1.0/ && npm start &
test-views-only:
	cd dashboard/v1.0/ && npm run lint
	cd dashboard/v1.0/ && npm run tlint
	cd dashboard/v1.0/ && prettier '**/*.tsx' --list-different
	cd dashboard/v1.0/ && npm run react-test
	cd dashboard/v1.0/ && npm run react-build
	cd dashboard/v1.0/ && npm run build
	cd dashboard/v1.0/ && npm start &

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

run-collector:
	go run src/collector/*.go

fix:
	go fmt ./...

lint:
	golangci-lint run

