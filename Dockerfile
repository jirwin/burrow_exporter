FROM golang:alpine as glide
RUN apk update
RUN apk add git
RUN go get github.com/Masterminds/glide
WORKDIR /go/src/github.com/jirwin/burrow_exporter
COPY . /go/src/github.com/jirwin/burrow_exporter
RUN glide install
RUN go build burrow-exporter.go

FROM alpine
COPY --from=glide /go/src/github.com/jirwin/burrow_exporter/burrow-exporter .
ENV BURROW_ADDR http://localhost:8000
ENV METRICS_ADDR 0.0.0.0:8080
ENV INTERVAL 30
ENV API_VERSION 2
CMD ./burrow-exporter --burrow-addr $BURROW_ADDR --metrics-addr $METRICS_ADDR --interval $INTERVAL --api-version $API_VERSION