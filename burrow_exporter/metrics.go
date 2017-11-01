package burrow_exporter

import "github.com/prometheus/client_golang/prometheus"

var (
	KafkaConsumerPartitionLag = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "kafka_burrow_partition_lag",
			Help: "The lag of the latest offset commit on a partition as reported by burrow.",
		},
		[]string{"cluster", "group", "topic", "partition"},
	)
	KafkaConsumerPartitionCurrentOffset = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "kafka_burrow_partition_current_offset",
			Help: "The latest offset commit on a partition as reported by burrow.",
		},
		[]string{"cluster", "group", "topic", "partition"},
	)
	KafkaConsumerPartitionMaxOffset = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "kafka_burrow_partition_max_offset",
			Help: "The log end offset on a partition as reported by burrow.",
		},
		[]string{"cluster", "group", "topic", "partition"},
	)
	KafkaConsumerTotalLag = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "kafka_burrow_total_lag",
			Help: "The total amount of lag for the consumer group as reported by burrow",
		},
		[]string{"cluster", "group"},
	)
)

func init() {
	prometheus.MustRegister(KafkaConsumerPartitionLag)
	prometheus.MustRegister(KafkaConsumerPartitionCurrentOffset)
	prometheus.MustRegister(KafkaConsumerPartitionMaxOffset)
	prometheus.MustRegister(KafkaConsumerTotalLag)
}
