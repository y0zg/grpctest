FROM golang:alpine as builder

COPY . /go/src/github.com/zenoss/grpctest
WORKDIR /go/src/github.com/zenoss/grpctest
RUN go build -o /bin/grpctest

FROM alpine
COPY --from=builder /bin/grpctest /bin/grpctest
ENTRYPOINT ["/bin/grpctest"]
