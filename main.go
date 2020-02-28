package main

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	"context"
	"os/signal"
	"syscall"

	"github.com/jirwin/burrow_exporter/exporter"
	"github.com/urfave/cli/v2"
)

var Version = "0.0.5"

func main() {
	app := cli.App{
		Name:    "burrow-exporter",
		Version: Version,
		Action: func(c *cli.Context) error {
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

			exporter := exporter.MakeBurrowExporter(
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
		},

		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "burrow-addr",
				Usage:   "Address that burrow is listening on",
				EnvVars: []string{"BURROW_ADDR"},
			},
			&cli.StringFlag{
				Name:    "metrics-addr",
				Usage:   "Address to run prometheus on",
				EnvVars: []string{"METRICS_ADDR"},
			},
			&cli.IntFlag{
				Name:    "interval",
				Usage:   "The interval(seconds) specifies how often to scrape burrow.",
				EnvVars: []string{"INTERVAL"},
			},
			&cli.IntFlag{
				Name:    "api-version",
				Usage:   "Burrow API version to leverage",
				Value:   2,
				EnvVars: []string{"API_VERSION"},
			},
			&cli.BoolFlag{
				Name:    "skip-partition-status",
				Usage:   "Skip exporting the per-partition status",
				EnvVars: []string{"SKIP_PARTITION_STATUS"},
			},
			&cli.BoolFlag{
				Name:    "skip-group-status",
				Usage:   "Skip exporting the per-group status",
				EnvVars: []string{"SKIP_GROUP_STATUS"},
			},
			&cli.BoolFlag{
				Name:    "skip-partition-lag",
				Usage:   "Skip exporting the partition lag",
				EnvVars: []string{"SKIP_PARTITION_LAG"},
			},
			&cli.BoolFlag{
				Name:    "skip-partition-current-offset",
				Usage:   "Skip exporting the current offset per partition",
				EnvVars: []string{"SKIP_PARTITION_CURRENT_OFFSET"},
			},
			&cli.BoolFlag{
				Name:    "skip-partition-max-offset",
				Usage:   "Skip exporting the partition max offset",
				EnvVars: []string{"SKIP_PARTITION_MAX_OFFSET"},
			},
			&cli.BoolFlag{
				Name:    "skip-total-lag",
				Usage:   "Skip exporting the total lag",
				EnvVars: []string{"SKIP_TOTAL_LAG"},
			},
			&cli.BoolFlag{
				Name:    "skip-topic-partition-offset",
				Usage:   "Skip exporting topic partition offset",
				EnvVars: []string{"SKIP_TOPIC_PARTITION_OFFSET"},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("error running burrow-exporter")
		os.Exit(1)
	}
}
