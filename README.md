# burrow-exporter

A simple prometheus exporter for gathering Kafka consumer group info
from [burrow](https://github.com/linkedin/Burrow).


## Docker usage

	docker build -t burrow_exporter --rm .
	docker run --rm burrow_exporter --interval 30 --burrow-addr http://${BURROW_IP}:8000
