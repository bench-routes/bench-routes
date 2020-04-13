# Bench-routes

[![Build Status](https://travis-ci.com/zairza-cetb/bench-routes.svg?branch=master)](https://travis-ci.com/zairza-cetb/bench-routes)
[![Go Report Card](https://goreportcard.com/badge/github.com/zairza-cetb/bench-routes)](https://goreportcard.com/report/github.com/zairza-cetb/bench-routes)
[![Gitter](https://img.shields.io/badge/join%20discussions%20on%20gitter-%23benchroutes-green/)](https://gitter.im/bench-routes/community#)

Bench-routes is a highly scalable network benchmarking, routes performance and monitoring tool, that monitors in regular intervals the state
of the server, running as a daemon process.

#### Dataflow

![Screenshot from 2020-03-21 20-09-00](https://user-images.githubusercontent.com/33792202/77228928-b139e900-6bb0-11ea-877b-54afffd2aa53.png)

**For more information, read the [docs](https://docs.google.com/document/d/1jGfc2eXvToRL9anzosTLQ4zJ7fdFxMGfaiDv2BYHEvw/edit?usp=sharing)**

### Read the complete idea and approach [here](https://github.com/zairza-cetb/bench-routes/blob/master/approach.md).

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
Bench Routes uses DCO(Developer Certificate Origin) to certify that the contributor wrote the particular code or otherwise have the right to submit the code they are contributing to the project.For complete details on DCO  <a href="https://probot.github.io/apps/dco/" target="_blank">Click Here</a>.

Follow the below `commit` syntax to certify the code and pass the DCO test.
```
git commit -s -m <commit-message>
```

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

## Postman Usage
1. Download [Postman](https://www.postman.com/downloads/) and Install it.
2. Create a new collection.

### To Check Service State
1. Add request
2. Select method **GET**
3. Copy and Enter below request url  
`http://localhost:9090/service-state` 
4. Send the request to url.
5. This API returns the state of the services (active or passive) in real-time.

### To Get Routes Summary
1. Add request
2. Select method **GET**
3. Copy and Enter below request url  
`http://localhost:9090/routes-summary` 
4. Send the request to url.
5. This API returns the list of all URLs/Routes that are being monitored for testing using the application.

For more information, regarding usage in different languages. Visit [Bench-Routes](https://documenter.getpostman.com/view/6521254/SzRuWqq9?version=latest).

## GSSoC'20
**Bench-routes** has been selected at [GirlScript Summer of Code 2020](https://www.gssoc.tech/)!!! -a three months long Open Source program conducted by [GirlScript Foundation](https://www.girlscript.tech/home).

### 👨 Project Owner

- Harkishen Singh <p>[<img src="https://img.icons8.com/windows/32/000000/github-2.png" display = "inline-block">](https://github.com/Harkishen-Singh) [<img src="https://img.icons8.com/ios-glyphs/30/000000/linkedin-2.png"/>](https://www.linkedin.com/in/harkishen-singh/)</p>

### 👬  Mentors

- Aquib Baig <p>[<img src="https://img.icons8.com/windows/32/000000/github-2.png" display = "inline-block">](https://github.com/aquibbaig) [<img src="https://img.icons8.com/ios-glyphs/30/000000/linkedin-2.png"/>](https://www.linkedin.com/in/baigaquib/)</p>

- Ganesh Patro <p>[<img src="https://img.icons8.com/windows/32/000000/github-2.png" display = "inline-block">](https://github.com/ganeshpatro321) [<img src="https://img.icons8.com/ios-glyphs/30/000000/linkedin-2.png"/>](https://www.linkedin.com/in/ganeshpatro321/)</p>

- Marut Pandya <p>[<img src="https://img.icons8.com/windows/32/000000/github-2.png" display = "inline-block">](https://github.com/pandyamarut) [<img src="https://img.icons8.com/ios-glyphs/30/000000/linkedin-2.png"/>](https://www.linkedin.com/in/marut-pandya-a16480131/)</p>

- Vatsal Trivedi <p>[<img src="https://img.icons8.com/windows/32/000000/github-2.png" display = "inline-block">](https://github.com/vattytrivedi) [<img src="https://img.icons8.com/ios-glyphs/30/000000/linkedin-2.png"/>](https://www.linkedin.com/in/trivedi-vatsal/)</p>

Feel free to ask your queries!! 🙌

### Slack Channel

- [#proj_bench-routes](https://app.slack.com/client/TRN1H1V43/CUCKG6R0V)



