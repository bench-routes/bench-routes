## Installation instructions 

Please follow the following steps in order to set up your development
environment.

### Docker

1. Make sure [docker](https://www.docker.com/) is installed, and your user is in the docker group.
2. Run `docker build -t bench-routes .`
3. Run `docker run -p 9090:9090 -it bench-routes`

### Local machine

1. Install Golang +1.10.x (assuming `git` is already installed).
2. Make a default repository for cloning the project. This should be strictly inside the `GOPATH`. Paste this instruction in your terminal to get started.
`mkdir -p $GOPATH/src/github.com/zairza-cetb/`.
3. Navigate to the directory via `cd $GOPATH/src/github.com/zairza-cetb`.
4. Clone the repository via `git clone https://github.com/zairza-cetb/bench-routes.git`.
5. Navigate into the cloned repo `cd bench-routes/`.
6. Install all dependencies via `go get -v -u ./...`.
7. To start running, `make run` will start the service. For running independent modules, make the .go files in the modules
as `package main` and include in them, a `main()` function. This is just for testing when developing independent application module. Make sure to link everything with `main.go` file or the parent file before pushing, else the **CI-builds** will fail.

### Optional
1. Installing [golangci-lint](https://github.com/golangci/golangci-lint) (simply paste command in your terminal after each step): 

```
curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s sh -s v1.18.0

golangci-lint --version

golangci-lint run
```
2. Scanning your `.go` files for linting issues using `golangci-lint`


Please feel free to open any issue in case you encounter any issues while setting up the development environment.
