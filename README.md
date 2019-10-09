# bench-routes

[![Build Status](https://travis-ci.com/zairza-cetb/bench-routes.svg?branch=master)](https://travis-ci.com/zairza-cetb/bench-routes)
[![Gitter](https://img.shields.io/badge/join%20discussions%20on%20gitter-%23benchroutes-green/)](https://gitter.im/bench-routes/community#)

bench-routes is a highly scalable network benchmarking, routes performance and monitoring tool, that monitors in regular intervals the state
of the server, running as a daemon process. For more information, read the [docs](https://docs.google.com/document/d/1jGfc2eXvToRL9anzosTLQ4zJ7fdFxMGfaiDv2BYHEvw/edit?usp=sharing)

Monitoring has been tough and with the increase in the routes used in any sophisticated project, the performance and metrics of an application are seriously affected.
With an increase in server computational models, the probability of a complete request-response cycle without any throws is nowhere close to 1. 

The primary goals of the project are:

```
1. Benchmark route
  (a) Load-handling of application on the individual route.
  (b) Test various possibilities of data in params (Permutate params), like sending an empty param
      to see how the server response behaves.
2. Analyse network performance of the hosted application irrespectively of containerization
  (a) Network ping
  (b) Jitter analysis
  (c) Packet loss
3. Log error handling capability of the application
4. Maintain a check on server-route output and alert on changes above the threshold
5. Graphical view using ElectronJS
```

For installation instructions, please head-over to [INSTALL.md](https://github.com/zairza-cetb/bench-routes/blob/master/INSTALL.md).

## Making Commits in bench-routes
Bench Routes uses DCO(Developer Certificate Origin) to certify that the contributor wrote the partcular code or otherwise have the right to submit the code they are contributing to the project.For complete details on DCO  <a href="https://probot.github.io/apps/dco/" target="_blank">Click Here</a>.

Follow the below `commit` syntax to cerify the code and pass the DCO test.
```
git commit -s -m <commit-message>
```

## Use of MakeFile in bench-routes
We use `make` for building and executing the program.

Follow the commands to make the development proces easier:

1. Updating the dependencies: `make update`
2. Executing the application (assuming all dependencies are installed): `make run`
2. Building the application for the current OS: `make build`
3. Testing Golang code: `make test`
4. Complete testing include building for all OSs out there: `make test_complete`
5. Cleaning up the residula files: `make clean`
6. *(optional)* Check linting (assuming [golangci-lint](https://github.com/golangci/golangci-lint#install) is installed): `make lint`
