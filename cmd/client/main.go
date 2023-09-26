package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/natalya-revtova/system-monitoring/pkg/client"
)

var host, port, notifyInterval, averageInterval string

func init() {
	flag.StringVar(&host, "host", "127.0.0.1", "Host of monitoring service to connect")
	flag.StringVar(&port, "port", "50051", "Port of monitoring service to connect")
	flag.StringVar(&notifyInterval, "notify-interval", "5", "Interval for results printing (sec)")
	flag.StringVar(&averageInterval, "average-interval", "10", "Interval for average calculation (sec)")
}

func main() {
	notifyInt, err := strconv.Atoi(notifyInterval)
	if err != nil {
		fmt.Printf("invalid notify-interval parameter: %v\n", err)
		os.Exit(1)
	}
	avgInt, err := strconv.Atoi(averageInterval)
	if err != nil {
		fmt.Printf("invalid average-interval parameter: %v\n", err)
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	cli := client.NewClient(host, port)
	if err := cli.Start(ctx, notifyInt, avgInt); err != nil {
		fmt.Printf("Start client error: %v\n", err)
		cancel()
	}
}
