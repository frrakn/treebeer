FROM golang

ADD . /go/src/github.com/frrakn/treebeer/context/contextManager

WORKDIR /go/src/github.com/frrakn/treebeer/context/contextManager

RUN go get github.com/go-sql-driver/mysql
RUN go get github.com/golang/glog
RUN go get github.com/golang/protobuf/proto
RUN go get github.com/jmoiron/sqlx
RUN go get github.com/juju/errors
RUN go get golang.org/x/net/context
RUN go get google.golang.org/grpc
RUN go build /go/src/github.com/frrakn/treebeer/context/contextManager/manager.go

ENTRYPOINT ./main /go/src/github.com/frrakn/treebeer/context/contextManager/conf/default.cfg

EXPOSE 8080
