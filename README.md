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

#### API_VERSION
Burrow API version to leverage (default: 2)

### Example

```sh
# build docker image
docker build -t burrow_exporter .

# with env variables
docker run -d -p 8080:8080 \
  -e BURROW_ADDR="http://localhost:8000" \
  -e METRICS_ADDR="0.0.0.0:8080" \
  -e INTERVAL="30" \
  -e API_VERSION="2" \
  burrow_exporter
# with custom command
docker run -d -p 8080:8080 burrow_exporter ./burrow-exporter --burrow-addr http://localhost:8000 --metrics-addr 0.0.0.0:8080 --interval 30 --api-version 2

```
