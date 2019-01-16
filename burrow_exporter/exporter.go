package burrow_exporter

import (
	"context"

	"sync"
	"time"

	"net/http"

	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type BurrowExporter struct {
	client                     *BurrowClient
	metricsListenAddr          string
	interval                   int
	wg                         sync.WaitGroup
	skipPartitionStatus        bool
	skipConsumerStatus         bool
	skipPartitionLag           bool
	skipPartitionCurrentOffset bool
	skipPartitionMaxOffset     bool
	skipTotalLag               bool
	skipTopicPartitionOffset   bool
}

func (be *BurrowExporter) processGroup(cluster, group string) {
	status, err := be.client.ConsumerGroupLag(cluster, group)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("error getting status for consumer group. returning.")
		return
	}

	for _, partition := range status.Status.Partitions {
		if !be.skipPartitionLag {
			KafkaConsumerPartitionLag.With(prometheus.Labels{
				"cluster":   status.Status.Cluster,
				"group":     status.Status.Group,
				"topic":     partition.Topic,
				"partition": strconv.Itoa(int(partition.Partition)),
			}).Set(float64(partition.CurrentLag))
		}

		if !be.skipPartitionCurrentOffset {
			KafkaConsumerPartitionCurrentOffset.With(prometheus.Labels{
				"cluster":   status.Status.Cluster,
				"group":     status.Status.Group,
				"topic":     partition.Topic,
				"partition": strconv.Itoa(int(partition.Partition)),
			}).Set(float64(partition.End.Offset))
		}

		if !be.skipPartitionStatus {
			KafkaConsumerPartitionCurrentStatus.With(prometheus.Labels{
				"cluster":   status.Status.Cluster,
				"group":     status.Status.Group,
				"topic":     partition.Topic,
				"partition": strconv.Itoa(int(partition.Partition)),
			}).Set(float64(Status[partition.Status]))
		}

		if !be.skipPartitionMaxOffset {
			KafkaConsumerPartitionMaxOffset.With(prometheus.Labels{
				"cluster":   status.Status.Cluster,
				"group":     status.Status.Group,
				"topic":     partition.Topic,
				"partition": strconv.Itoa(int(partition.Partition)),
			}).Set(float64(partition.End.MaxOffset))
		}
	}

	if !be.skipTotalLag {
		KafkaConsumerTotalLag.With(prometheus.Labels{
			"cluster": status.Status.Cluster,
			"group":   status.Status.Group,
		}).Set(float64(status.Status.TotalLag))
	}

	if !be.skipConsumerStatus {
		KafkaConsumerStatus.With(prometheus.Labels{
			"cluster": status.Status.Cluster,
			"group":   status.Status.Group,
		}).Set(float64(Status[status.Status.Status]))
	}
}

func (be *BurrowExporter) processTopic(cluster, topic string) {
	details, err := be.client.ClusterTopicDetails(cluster, topic)
	if err != nil {
		log.WithFields(log.Fields{
			"err":   err,
			"topic": topic,
		}).Error("error getting status for cluster topic. returning.")
		return
	}

	if !be.skipTopicPartitionOffset {
		for i, offset := range details.Offsets {
			KafkaTopicPartitionOffset.With(prometheus.Labels{
				"cluster":   cluster,
				"topic":     topic,
				"partition": strconv.Itoa(i),
			}).Set(float64(offset))
		}
	}
}

func (be *BurrowExporter) processCluster(cluster string) {
	groups, err := be.client.ListConsumers(cluster)
	if err != nil {
		log.WithFields(log.Fields{
			"err":     err,
			"cluster": cluster,
		}).Error("error listing consumer groups. returning.")
		return
	}

	topics, err := be.client.ListClusterTopics(cluster)
	if err != nil {
		log.WithFields(log.Fields{
			"err":     err,
			"cluster": cluster,
		}).Error("error listing cluster topics. returning.")
		return
	}

	wg := sync.WaitGroup{}

	for _, group := range groups.ConsumerGroups {
		wg.Add(1)

		go func(g string) {
			defer wg.Done()
			be.processGroup(cluster, g)
		}(group)
	}

	for _, topic := range topics.Topics {
		wg.Add(1)

		go func(t string) {
			defer wg.Done()
			be.processTopic(cluster, t)
		}(topic)
	}

	wg.Wait()
}

func (be *BurrowExporter) startPrometheus() {
	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(be.metricsListenAddr, nil)
}

func (be *BurrowExporter) Close() {
	be.wg.Wait()
}

func (be *BurrowExporter) Start(ctx context.Context) {
	be.startPrometheus()

	be.wg.Add(1)
	defer be.wg.Done()

	be.mainLoop(ctx)
}

func (be *BurrowExporter) scrape() {
	start := time.Now()
	log.WithField("timestamp", start.UnixNano()).Info("Scraping burrow...")
	clusters, err := be.client.ListClusters()
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("error listing clusters. Continuing.")
		return
	}

	wg := sync.WaitGroup{}

	for _, cluster := range clusters.Clusters {
		wg.Add(1)

		go func(c string) {
			defer wg.Done()
			be.processCluster(c)
		}(cluster)
	}

	wg.Wait()

	end := time.Now()
	log.WithFields(log.Fields{
		"timestamp": end.UnixNano(),
		"took":      end.Sub(start),
	}).Info("Finished scraping burrow.")
}

func (be *BurrowExporter) mainLoop(ctx context.Context) {
	timer := time.NewTicker(time.Duration(be.interval) * time.Second)

	// scrape at app start without waiting for the first interval to elapse
	be.scrape()

	for {
		select {
		case <-ctx.Done():
			log.Info("Shutting down exporter.")
			timer.Stop()
			return

		case <-timer.C:
			be.scrape()
		}
	}
}

func MakeBurrowExporter(burrowUrl string, apiVersion int, metricsAddr string, interval int, skipPartitionStatus bool,
	skipConsumerStatus bool, skipPartitionLag bool, skipPartitionCurrentOffset bool, skipPartitionMaxOffset bool, skipTotalLag bool, skipTopicPartitionOffset bool) *BurrowExporter {
	return &BurrowExporter{
		client:                     MakeBurrowClient(burrowUrl, apiVersion),
		metricsListenAddr:          metricsAddr,
		interval:                   interval,
		skipPartitionStatus:        skipPartitionStatus,
		skipConsumerStatus:         skipConsumerStatus,
		skipPartitionLag:           skipPartitionLag,
		skipPartitionCurrentOffset: skipPartitionCurrentOffset,
		skipPartitionMaxOffset:     skipPartitionMaxOffset,
		skipTotalLag:               skipTotalLag,
		skipTopicPartitionOffset:   skipTopicPartitionOffset,
	}
}
