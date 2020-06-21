update:
	echo "updating dependencies ..."
	go get -u ./...

build:
	echo "building bench-routes ..."
	go build src/main.go
	mv main bench-routes

view-v1.0:
	cd dashboard/v1.0/ && sudo npm start

view-v1.1:
	cd dashboard/v1.1/ && sudo yarn start

test-views-v1.0:
	cd dashboard/v1.0/ && npm install
	cd dashboard/v1.0/ && npm run lint
	cd dashboard/v1.0/ && npm run tlint
	cd dashboard/v1.0/ && npm run react-test
	cd dashboard/v1.0/ && npm run react-build
	cd dashboard/v1.0/ && npm run build
	cd dashboard/v1.0/ && npm start &

test-views-v1.1:
	cd dashboard/v1.1/ && yarn install
	cd dashboard/v1.1/ && yarn run lint
	cd dashboard/v1.1/ && yarn run tlint
	cd dashboard/v1.1/ && prettier '**/*.tsx' --list-different
	cd dashboard/v1.1/ && yarn run build
	cd dashboard/v1.1/ && yarn start &

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

build-frontend:
	cd dashboard/v1.1/ && yarn install && yarn build

test-non-verbose: build
	go clean -testcache
	go test ./...

test-services: build
	./bench-routes &
	cd tests && yarn install
	yarn global add mocha
	mocha tests/browser.js

test_complete: build
	./shell/go-build-all.sh
	echo "test success!"

run:
	echo "compiling go-code and executing bench-routes"
	echo "using 9090 as default service listener port"
	go run src/*.go 9090

run-collector:
	go run src/collector/main.go

fix:
	go fmt ./...
	cd dashboard/v1.1/ && npm run prettier-fix
	cd dashboard/v1.1/ && npm run tlint-fix

lint:
	golangci-lint run

