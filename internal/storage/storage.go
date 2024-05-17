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

func (s *Storage) Get(name string, n int) ([]models.Metrics, bool) {
	s.mux.Lock()
	defer s.mux.Unlock()

	if n <= 0 {
		return nil, false
	}

	if _, ok := s.metrics[name]; !ok {
		return nil, false
	}

	mlen := len(s.metrics[name])
	if mlen < n {
		return nil, true
	}
	return s.metrics[name][mlen-n:], true
}
