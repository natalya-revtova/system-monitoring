package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/natalya-revtova/system-monitoring/internal/config"
	"github.com/natalya-revtova/system-monitoring/internal/grabber"
	"github.com/natalya-revtova/system-monitoring/internal/logger"
	"github.com/natalya-revtova/system-monitoring/internal/monitoring"
	"github.com/natalya-revtova/system-monitoring/internal/server"
	"github.com/natalya-revtova/system-monitoring/internal/storage"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "../../configs/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	config, err := config.NewConfig(configFile)
	if err != nil {
		fmt.Printf("failed to read configuration file: %v\n", err)
		os.Exit(1)
	}

	log := logger.New(config.Logger.Level)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	grabber := grabber.NewGrabber(grabber.GetOptions(config.Metrics), log)
	storage := storage.NewStorage()

	svc := monitoring.NewService(ctx, grabber, storage, log)
	server := server.NewServer(config.Server, svc, log)

	go func() {
		<-ctx.Done()
		server.Stop()
	}()

	log.Info("Monitoring daemon is running")

	if err := server.Start(); err != nil {
		log.Error("Start server", "error", err)
		cancel()
	}
}
