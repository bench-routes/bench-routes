# Installation instructions 

Please follow the following steps in order to set up your development
environment.

## Bare metal (local)

1. Install Golang +1.10.x (assuming `git` is already installed).
2. Make a default repository for cloning the project. This should be strictly inside the `GOPATH`. Paste this instruction in your terminal to get started.
`mkdir -p $GOPATH/src/github.com/bench-routes`.
3. Navigate to the directory via `cd $GOPATH/src/github.com/bench-routes`.
4. Clone the repository via `git clone https://github.com/bench-routes/bench-routes.git`.
5. Navigate into the cloned repo `cd bench-routes`.
6. Install all dependencies via `go get -v -u ./...`.
7. To start running, `make run` will start the service.

## Docker

1. Make sure [docker](https://www.docker.com/) is installed, and your user is in the docker group.
2. Run `docker build -t bench-routes .`
3. Run `docker run -p 9990:9990 -it bench-routes`

Please feel free to open any [new-issue](https://github.com/bench-routes/bench-routes/issues/new/choose) in case you encounter with any issues while setting up the development environment.
