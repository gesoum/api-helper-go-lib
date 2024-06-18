FROM golang:1.22

WORKDIR /go/src/api-helper

COPY Makefile go.* ./

RUN make tools