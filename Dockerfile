FROM golang:1.14.2-alpine3.11
RUN apk add --no-cache bash
WORKDIR $GOPATH/src/github.com/zairza-cetb/bench-routes
COPY . $GOPATH/src/github.com/zairza-cetb/bench-routes
RUN rm -R $GOPATH/src/github.com/zairza-cetb/bench-routes/storage
RUN cd $GOPATH/src/github.com/zairza-cetb/bench-routes/src && go get -v ./...
RUN go build $GOPATH/src/github.com/zairza-cetb/bench-routes/src/main.go
RUN mv main bench-routes
EXPOSE 9090
CMD [ "./bench-routes" ]
