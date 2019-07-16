FROM golang:1.12 AS builder
ADD . /go/src/github.com//govargo/bar-controller
WORKDIR /go/src/github.com/govargo/bar-controller
RUN CGO_ENABLED=0 GO111MODULE=on go build -o bar-controller .

FROM alpine:3.10.1
WORKDIR /
COPY --from=builder /go/src/github.com/govargo/bar-controller .
CMD ["/bar-controller"]
