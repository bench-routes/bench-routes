# Installation Instructions 

Please Follow the following steps in order to set up your development
environment.

## Bare metal (Local)

1. Install Golang +1.10.x (assuming `git` is already installed).
```console
sudo add-apt-repository ppa:longsleep/golang-backports
sudo apt update
sudo apt install golang-go
```
2. Set up `$GOPATH` environment variable as your preferred directory.Run:
```bash
nano ~/.bash_profile
```
Enter `export GOPATH="directory Name"` and save the file. 

3. Make a default repository for cloning the project. This should be strictly inside the `GOPATH`. 
Paste this instruction in your terminal to get started.
```bash
mkdir -p $GOPATH/src/github.com/bench-routes
```
4. Navigate to the directory via 
```bash
cd $GOPATH/src/github.com/bench-routes
```
5. Clone the repository via 
```bash
git clone https://github.com/bench-routes/bench-routes.git
```
6. Navigate into the cloned repo 
```bash
cd bench-routes
```
7. Install all dependencies via 
```go
go get -v -u ./...
```
8. To start running the service. 
```bash 
make run
``` 

9. To setup and run the UI for the project kindly follow the instructions mentioned [here](https://github.com/bench-routes/dashboard#readme).

## Docker

1. Make sure [docker](https://www.docker.com/) is installed, and your user is in the docker group.
2. Run `docker build -t bench-routes .`
3. Run `docker run -p 9990:9990 -it bench-routes`

## Installation Instructions for Windows system using WSL

1. Install [WSL](https://docs.microsoft.com/en-us/windows/wsl/install-win10) in your windows machine.
```bash
 wsl install -d ubuntu
 ```
2. Install Golang in your WSL distro.
3. Set up `$GOPATH` environment variable as your preferred directory in WSL.Run:
```bash
nano ~/.bash_profile
```
Enter `export GOPATH="directory Name"` and save the file.

4. Open the project in VScode WSL window.
5. Install all dependencies via 
```go
go get -v -u ./...
```
6. To start running the service. 
```bash 
make run
``` 

Please feel free to open any [new-issue](https://github.com/bench-routes/bench-routes/issues/new/choose) in case you encounter with any issues while setting up the development environment.
