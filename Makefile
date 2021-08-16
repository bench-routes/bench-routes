## Builds the application for the current OS.
build:
	echo "building bench-routes ..."
	go build -o bench-routes src/main.go

## Removes all residual files.
clean:
	rm -R bench-routes

## Runs Golang unit tests
test:
	go clean -testcache
	go test -race ./...

## Executes the application (assuming all dependencies are installed)
run:
	echo "compiling go-code and executing bench-routes"
	echo "using 9990 as default service listener port"
	go run src/*.go 9990
