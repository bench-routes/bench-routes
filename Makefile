update:
	echo "updating dependencies ..."
	go get -u ./...

build:
	echo "building bench-routes ..." 
	go build src/main.go
	mv main bench-routes

clean:
	rm -R build/ bench-routes

test: build
	go clean -testcache
	go test -v ./...

test_complete: build
	./shell/go-build-all.sh
	echo "test success! cleaning ..."
	make clean

run:
	echo "compiling go-code and executing bench-routes"
	echo "using 9090 as default service listerner port"
	go run src/main.go 9090

fix:
	go fmt ./...

lint:
	golangci-lint run

