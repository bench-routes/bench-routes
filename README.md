# Bench-routes

![Linux build](https://github.com/bench-routes/bench-routes/actions/workflows/build-linux.yml/badge.svg)
![Mac build](https://github.com/bench-routes/bench-routes/actions/workflows/build-mac.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/bench-routes/bench-routes)](https://goreportcard.com/report/github.com/bench-routes/bench-routes)
[![Gitter](https://img.shields.io/badge/join%20discussions%20on%20gitter-%23benchroutes-green/)](https://gitter.im/bench-routes/community#)

Modern web applications can have routes ranging from a few to millions in numbers. This makes it tough to discover the
condition and state of such application at any given point. Bench-routes monitors the routes of a web application
and helps you know about the current state of each route, along with various related performance metrics.

### Dataflow

![Screenshot from 2020-03-21 20-09-00](https://user-images.githubusercontent.com/33792202/77228928-b139e900-6bb0-11ea-877b-54afffd2aa53.png)

### Primary goals
1. Monitoring web applications routes at scale.
2. Querying the monitored data in an interactive UI that is minimalistic to learn.
3. Reporting in case of any abnormalities.

For installation instructions, please head-over to [INSTALL.md](https://github.com/bench-routes/bench-routes/blob/master/INSTALL.md).

## Use of MakeFile in bench-routes
We use `make` for building and executing the program.

Follow the commands to make the development process easier:

1. Updating the dependencies: `make update`
2. Executing the application (assuming all dependencies are installed): `make run`
2. Building the application for the current OS: `make build`
3. Testing Golang code: `make test`
4. Complete testing include building for all OSs out there: `make test_complete`
5. Cleaning up the residual files: `make clean`
6. *(optional)* Check linting (assuming [golangci-lint](https://github.com/golangci/golangci-lint#install) is installed): `make lint`

## UI of Bench-routes
The old UI for the project is depreceated for now and no longer mantained. Kindly have a look at our new UI [here](https://github.com/bench-routes/dashboard).

## Postman Usage
1. Download [Postman](https://www.postman.com/downloads/) and Install it.
2. Create a new collection.

### To Check Service State
1. Add request
2. Select method **GET**
3. Copy and Enter below request url  
`http://localhost:9990/service-state` 
4. Send the request to url.
5. This API returns the state of the services (active or passive) in real-time.

### To Get Routes Summary
1. Add request
2. Select method **GET**
3. Copy and Enter below request url  
`http://localhost:9990/routes-summary` 
4. Send the request to url.
5. This API returns the list of all URLs/Routes that are being monitored for testing using the application.

For more information, regarding the APIs. Visit [Bench-Routes](https://documenter.getpostman.com/view/6521254/SzRuWqq9?version=latest).

### Communication

- Instant messaging: [bench-routes/community](https://gitter.im/bench-routes/community)
- Discussions (Recent): https://groups.google.com/forum/#!forum/bench-routes-discussion

## Want to contribute? [![contributions welcome](https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat)](https://github.com/bench-routes/bench-routes/issues)

Want to file a bug, request a feature, contribute some code, or improve documentation? Excellent! Read up on our guidelines for [contributing](CONTRIBUTING.md) and then check out one of our [issues](https://github.com/bench-routes/bench-routes/issues). Make sure you follow the guidelines before sending a contribution!
   
### Supported by

[![JetBrains](https://user-images.githubusercontent.com/36980003/111029750-5b46fc80-8424-11eb-800f-ad7f64844f0d.png)](https://jb.gg/OpenSource)
