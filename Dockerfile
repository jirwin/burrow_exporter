FROM golang:alpine

# ENV GO111MODULE=on
# ENV GOPROXY=https://goproxy.cn

WORKDIR /opt

COPY . .

RUN CGO_ENABLED=0 go build -ldflags "-s" -o burrow-exporter

FROM scratch

ENV BURROW_ADDR http://localhost:8000
ENV METRICS_ADDR 0.0.0.0:8080
ENV API_VERSION 3

COPY --from=0 /opt/burrow-exporter .

ENTRYPOINT ["/burrow-exporter"]
CMD "--burrow-addr" "$BURROW_ADDR" "--metrics-addr" "$METRICS_ADDR" "--interval" "$INTERVAL" "--api-version" "$API_VERSION"