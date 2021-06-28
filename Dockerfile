FROM golang:latest 

RUN mkdir /app

RUN echo $GOPATH

RUN mkdir -p $GOPATH/src/github.com/bhaskarhc

RUN mkdir -p $GOPATH/src/github.com/bhaskarhc/kube-infra-status

ADD . $GOPATH/src/github.com/bhaskarhc/ci-e2e-status

WORKDIR $GOPATH/src/github.com/bhaskarhc/ci-e2e-status

# RUN go mod init

RUN go mod vendor

RUN go build -o /app/main .

CMD ["/app/main"]

EXPOSE 3000