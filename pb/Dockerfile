FROM golang:latest as builder


RUN apt update && apt install unzip
RUN cd / && wget https://github.com/protocolbuffers/protobuf/releases/download/v3.7.1/protoc-3.7.1-linux-x86_64.zip
RUN cd / && unzip protoc* -d protoc3
RUN mv /protoc3/bin/* /usr/local/bin/
RUN mv /protoc3/include/* /usr/local/include/

RUN go get -u github.com/golang/protobuf/protoc-gen-go
ENV SRC_DIR /go/src/github.com/zenoss/grpctest/pb
ENV GOOGLEAPIS_DIR /go/src/github.com/zenoss/grpctest/googleapis
WORKDIR /go/src/github.com/zenoss/grpctest

#protoc --include_imports --include_source_info --proto_path=${GOOGLEAPIS_DIR} --proto_path=$(pwd) --descriptor_set_out=pb/api_descriptor.pb pb/grpc_test.proto


#protoc --proto_path=$(pwd)/googleapis -I=$(pwd) --go_out=$(pwd) $(pwd)/pb/grpc_test.proto
