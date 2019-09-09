### Installation instructions 

Please follow the following steps in order to set up your development
environment.

1. Install Golang +1.10.x (assuming `git` is already installed).
2. Make a default repository for cloing the project. This should be strictly inside the `GOPATH`. Paste this instruction in your terminal to get started.
`mkdir -p $GOPATH/src/github.com/zairza-cetb/`.
3. Navigate to the directory via `cd $GOPATH/src/github.com/zairza-cetb`.
4. Clone the repository via `git clone https://github.com/zairza-cetb/bench-routes.git`.
5. Navigate into the cloned repo `cd bench-routes/`.
6. Install all dependencies via `go get -v -u ./...`.
7. To start running, `go run src/main.go 9090` will start the service. For running independent modules, make the .go files in the modules
as `package main` and include in them, a `main()` function. This is just for testing when developing independent application module. Make sure to link everything with `main.go` file or the parent file before pushing, else the **CI-builds** will fail.


Please feel free to open any issue in case you encounter any issues while setting up the development environment.
