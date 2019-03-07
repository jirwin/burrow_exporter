package main

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"

	"context"
	"os/signal"
	"syscall"

	"github.com/jirwin/burrow_exporter/burrow_exporter"
)

var Version = "0.0.5"

func main() {
	app := cli.NewApp()
	app.Version = Version
	app.Name = "burrow-exporter"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "burrow-addr",
			Usage:  "Address that burrow is listening on",
			EnvVar: "BURROW_ADDR",
		},
		cli.StringFlag{
			Name:   "metrics-addr",
			Usage:  "Address to run prometheus on",
			EnvVar: "METRICS_ADDR",
		},
		cli.IntFlag{
			Name:   "interval",
			Usage:  "The interval(seconds) specifies how often to scrape burrow.",
			EnvVar: "INTERVAL",
		},
		cli.IntFlag{
			Name:   "api-version",
			Usage:  "Burrow API version to leverage",
			Value:  2,
			EnvVar: "API_VERSION",
		},
		cli.BoolFlag{
			Name:   "skip-partition-status",
			Usage:  "Skip exporting the per-partition status",
			EnvVar: "SKIP_PARTITION_STATUS",
		},
		cli.BoolFlag{
			Name:   "skip-group-status",
			Usage:  "Skip exporting the per-group status",
			EnvVar: "SKIP_GROUP_STATUS",
		},
		cli.BoolFlag{
			Name:   "skip-partition-lag",
			Usage:  "Skip exporting the partition lag",
			EnvVar: "SKIP_PARTITION_LAG",
		},
		cli.BoolFlag{
			Name:   "skip-partition-current-offset",
			Usage:  "Skip exporting the current offset per partition",
			EnvVar: "SKIP_PARTITION_CURRENT_OFFSET",
		},
		cli.BoolFlag{
			Name:   "skip-partition-max-offset",
			Usage:  "Skip exporting the partition max offset",
			EnvVar: "SKIP_PARTITION_MAX_OFFSET",
		},
		cli.BoolFlag{
			Name:   "skip-total-lag",
			Usage:  "Skip exporting the total lag",
			EnvVar: "SKIP_TOTAL_LAG",
		},
		cli.BoolFlag{
			Name:   "skip-topic-partition-offset",
			Usage:  "Skip exporting topic partition offset",
			EnvVar: "SKIP_TOPIC_PARTITION_OFFSET",
		},
		cli.IntFlag{
			Name:   "verbosity",
			Usage:  "Set the log level",
			EnvVar: "LOG_LEVEL",
		},
	}

	app.Action = func(c *cli.Context) error {
		if !c.IsSet("burrow-addr") {
			fmt.Println("A burrow address is required (e.g. --burrow-addr http://localhost:8000)")
			os.Exit(1)
		}

		if !c.IsSet("metrics-addr") {
			fmt.Println("An address to run prometheus on is required (e.g. --metrics-addr localhost:8080)")
			os.Exit(1)
		}

		if !c.IsSet("interval") {
			fmt.Println("A scrape interval is required (e.g. --interval 30)")
			os.Exit(1)
		}

		done := make(chan os.Signal, 1)

		signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

		ctx, cancel := context.WithCancel(context.Background())

		burrow_exporter.SetLogLevel(c.Int("verbosity"))

		exporter := burrow_exporter.MakeBurrowExporter(
			c.String("burrow-addr"),
			c.Int("api-version"),
			c.String("metrics-addr"),
			c.Int("interval"),
			c.Bool("skip-partition-status"),
			c.Bool("skip-group-status"),
			c.Bool("skip-partition-lag"),
			c.Bool("skip-partition-current-offset"),
			c.Bool("skip-partition-max-offset"),
			c.Bool("skip-total-lag"),
			c.Bool("skip-topic-partition-offset"))
		go exporter.Start(ctx)

		<-done
		cancel()

		exporter.Close()

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("error running burrow-exporter")
		os.Exit(1)
	}
}
