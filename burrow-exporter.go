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

var Version = "0.0.4"

func main() {
	app := cli.NewApp()
	app.Version = Version
	app.Name = "burrow-exporter"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "burrow-addr",
			Usage: "Address that burrow is listening on",
		},
		cli.StringFlag{
			Name:  "metrics-addr",
			Usage: "Address to run prometheus on",
		},
		cli.IntFlag{
			Name:  "interval",
			Usage: "The interval(seconds) specifies how often to scrape burrow.",
		},
		cli.IntFlag{
			Name:  "api-version",
			Usage: "Burrow API version to leverage",
			Value: 2,
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

		exporter := burrow_exporter.MakeBurrowExporter(c.String("burrow-addr"), c.Int("api-version"), c.String("metrics-addr"), c.Int("interval"))
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
