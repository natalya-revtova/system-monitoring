package monitoring

import (
	"context"
	"math"
	"sync"
	"time"

	"github.com/natalya-revtova/system-monitoring/internal/logger"
	"github.com/natalya-revtova/system-monitoring/internal/models"
)

//go:generate go run github.com/vektra/mockery/v2@v2.33.0 --name=Grabber
type Grabber interface {
	Grab(results chan models.Metrics)
}

//go:generate go run github.com/vektra/mockery/v2@v2.33.0 --name=Storage
type Storage interface {
	Get(name string, n int) ([]models.Metrics, bool)
	Save(metrics models.Metrics)
}

type Service struct {
	grabber Grabber
	storage Storage
	options []string
	log     logger.ILogger
}

func NewService(ctx context.Context, grabber Grabber, storage Storage, options []string, log logger.ILogger) *Service {
	svc := Service{
		grabber: grabber,
		storage: storage,
		options: options,
		log:     log,
	}

	go svc.runCollector(ctx)
	return &svc
}

func (s *Service) runCollector(ctx context.Context) {
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			s.log.Info("Stop metrics collector")
			return
		case <-ticker.C:
			s.collect(make(chan models.Metrics, 5))
		}
	}
}

func (s *Service) collect(results chan models.Metrics) {
	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		s.grabber.Grab(results)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	for metrics := range results {
		s.storage.Save(metrics)
	}
}

func (s *Service) MetricsSnapshot(n int) []models.Metrics {
	result := make([]models.Metrics, len(s.options))

	for i, option := range s.options {
		metrics, ok := s.storage.Get(option, n)
		if !ok {
			continue
		}

		for len(metrics) == 0 {
			time.Sleep(time.Second)
			metrics, _ = s.storage.Get(option, n)
		}
		result[i] = calculateAverage(metrics)
	}
	return result
}

func calculateAverage(metrics []models.Metrics) models.Metrics {
	groupAvg := make([]models.Group, len(metrics[0].Groups))
	for i := range groupAvg {
		groupAvg[i] = models.Group{
			Metrics: make([]models.Metric, len(metrics[0].Groups[0].Metrics)),
		}
	}

	for i := range metrics {
		for j := range groupAvg {
			for k := range groupAvg[j].Metrics {
				groupAvg[j].Metrics[k].Value += metrics[i].Groups[j].Metrics[k].Value
			}
			groupAvg[j].Labels = metrics[i].Groups[j].Labels
		}
	}

	totalLen := len(metrics)
	for i := range groupAvg {
		for j, sum := range groupAvg[i].Metrics {
			groupAvg[i].Metrics[j].Value = math.Round(100*sum.Value/float64(totalLen)) / 100
			groupAvg[i].Metrics[j].Name = metrics[0].Groups[i].Metrics[j].Name
		}
	}

	return models.Metrics{
		Name:   metrics[0].Name,
		Groups: groupAvg,
	}
}
