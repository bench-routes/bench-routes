FROM golang:1.14.2-alpine3.11
RUN apk add --no-cache bash nodejs npm
RUN node --version
RUN npm --version
WORKDIR $GOPATH/src/github.com/bench-routes/bench-routes
COPY . $GOPATH/src/github.com/bench-routes/bench-routes
RUN cd $GOPATH/src/github.com/bench-routes/bench-routes/src && go get -v ./...
RUN go build $GOPATH/src/github.com/bench-routes/bench-routes/src/main.go
RUN mv main bench-routes
RUN npm install -g react-scripts yarn
RUN cd dashboard/v1.1/ && yarn install && yarn run build
RUN rm -R ui-builds/v1.1
RUN cp -r dashboard/v1.1/build ui-builds/v1.1
EXPOSE 9090
CMD [ "./bench-routes" ]
