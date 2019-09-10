init:
	echo "running make"

update:
	echo "updating dependencies ..."
	go get -u ./...

build: update
	echo "building bench-routes ..." 
	go build src/main.go
	mv main bench-routes

clean:
	rm -R build/ bench-routes

test: build
	go test ./...

test_complete: build
	./shell/go-build-all.sh
	echo "test success! cleaning ..."
	clean


