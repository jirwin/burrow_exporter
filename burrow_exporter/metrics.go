package burrow_exporter

import "github.com/prometheus/client_golang/prometheus"

var (
	KafkaLag = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "kafka_burrow_consumer_lag",
			Help: "The lag of the latest offset commit on a partition as reported by burrow.",
		},
		[]string{"cluster", "group", "topic", "partition"},
	)
)

func init() {
	prometheus.MustRegister(KafkaLag)
}
