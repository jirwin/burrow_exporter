FROM golang:alpine as build
WORKDIR /go/src/app
COPY . /go/src/app
RUN apk add --no-cache git && go-wrapper download && go-wrapper install

FROM alpine
COPY --from=build /go/bin/app /burrow_exporter
ENTRYPOINT ["/burrow_exporter", "--metrics-addr" ,"0.0.0.0:8080"]
