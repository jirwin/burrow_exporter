# burrow-exporter

A simple prometheus exporter for gathering Kafka consumer group info
from [burrow](https://github.com/linkedin/Burrow).


## Run with Docker

### required environment variables
#### BURROW_ADDR
A burrow address is required. Default: http://localhost:8000
#### METRICS_ADDR
An address to run prometheus on is required. Default: 0.0.0.0:8080
#### INTERVAL
A scrape interval is required. Default: 30

### Example
```sh
# with env variables
docker run \
  -e BURROW_ADDR="http://localhost:8000" \
  -e METRICS_ADDR="0.0.0.0:8080" \
  -e INTERVAL="30" \
  saada/burrow_exporter
# with custom command
docker run -d saada/burrow_exporter ./burrow-exporter --burrow-addr http://localhost:8000 --metrics-addr 0.0.0.0:8080 --interval 30
```