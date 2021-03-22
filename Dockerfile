FROM centos:centos8
RUN yum -y upgrade
RUN yum install -y net-tools nodejs npm golang
RUN node --version
RUN npm --version
WORKDIR $GOPATH/src/github.com/bench-routes/bench-routes
COPY . $GOPATH/src/github.com/bench-routes/bench-routes
RUN cd $GOPATH/src/github.com/bench-routes/bench-routes/src && go get -v ./...
RUN go build $GOPATH/src/github.com/bench-routes/bench-routes/src/main.go
RUN mv main bench-routes
RUN npm install -g yarn
RUN cd dashboard/v1.1/ && yarn install && yarn run build
RUN rm -R ui-builds/v1.1
RUN cp -r dashboard/v1.1/build ui-builds/v1.1
EXPOSE 9990
CMD [ "./bench-routes" ]
