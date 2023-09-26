package storage

import (
	"sync"

	"github.com/natalya-revtova/system-monitoring/internal/models"
)

var buffer = 60 // 36000 // store metrics for last 10 hours

type Storage struct {
	metrics map[string][]models.Metrics
	mux     sync.Mutex
}

func NewStorage() *Storage {
	return &Storage{
		metrics: make(map[string][]models.Metrics),
	}
}

func (s *Storage) Save(metrics models.Metrics) {
	s.mux.Lock()
	defer s.mux.Unlock()

	if _, ok := s.metrics[metrics.Name]; !ok {
		s.metrics[metrics.Name] = make([]models.Metrics, 0, buffer)
	}
	if len(s.metrics[metrics.Name]) == buffer {
		s.metrics[metrics.Name] = s.metrics[metrics.Name][1:]
	}
	s.metrics[metrics.Name] = append(s.metrics[metrics.Name], metrics)
}

func (s *Storage) Get(name string, n int) []models.Metrics {
	s.mux.Lock()
	defer s.mux.Unlock()

	mlen := len(s.metrics[name])
	if mlen < n || n <= 0 {
		return nil
	}
	return s.metrics[name][mlen-n:]
}
