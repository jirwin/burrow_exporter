FROM golang:1.14.0-alpine as builder
WORKDIR /go/src/github.com/jirwin/burrow_exporter
COPY . /go/src/github.com/jirwin/burrow_exporter
RUN go build

FROM alpine:3.9.5
COPY --from=builder /go/src/github.com/jirwin/burrow_exporter/burrow_exporter .
ENV BURROW_ADDR http://localhost:8000
ENV METRICS_ADDR 0.0.0.0:8080
ENV INTERVAL 30
ENV API_VERSION 2
CMD ./burrow-exporter --burrow-addr $BURROW_ADDR --metrics-addr $METRICS_ADDR --interval $INTERVAL --api-version $API_VERSION
