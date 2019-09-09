# bench-routes

[![Build Status](https://travis-ci.com/zairza-cetb/bench-routes.svg?branch=master)](https://travis-ci.com/zairza-cetb/bench-routes)


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
