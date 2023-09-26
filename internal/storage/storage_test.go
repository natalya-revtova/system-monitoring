package storage

import (
	"testing"

	"github.com/natalya-revtova/system-monitoring/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGet(t *testing.T) {
	cases := []struct {
		title    string
		metrics  []models.Metrics
		expected []models.Metrics
		name     string
		n        int
	}{
		{
			title: "get last 3 values of load_avg",
			n:     3,
			name:  "load_avg",
			metrics: []models.Metrics{
				{Name: "load_avg", Groups: []models.Group{{Metrics: []models.Metric{{Name: "avg1", Value: 1}}}}},
				{Name: "load_avg", Groups: []models.Group{{Metrics: []models.Metric{{Name: "avg1", Value: 2}}}}},
				{Name: "load_avg", Groups: []models.Group{{Metrics: []models.Metric{{Name: "avg1", Value: 3}}}}},
				{Name: "load_avg", Groups: []models.Group{{Metrics: []models.Metric{{Name: "avg1", Value: 2}}}}},
				{Name: "load_avg", Groups: []models.Group{{Metrics: []models.Metric{{Name: "avg1", Value: 1}}}}},
			},
			expected: []models.Metrics{
				{Name: "load_avg", Groups: []models.Group{{Metrics: []models.Metric{{Name: "avg1", Value: 3}}}}},
				{Name: "load_avg", Groups: []models.Group{{Metrics: []models.Metric{{Name: "avg1", Value: 2}}}}},
				{Name: "load_avg", Groups: []models.Group{{Metrics: []models.Metric{{Name: "avg1", Value: 1}}}}},
			},
		},
		{
			title: "get more values that is present in storage",
			n:     5,
			name:  "load_avg",
			metrics: []models.Metrics{
				{Name: "load_avg", Groups: []models.Group{{Metrics: []models.Metric{{Name: "avg1", Value: 1}}}}},
				{Name: "load_avg", Groups: []models.Group{{Metrics: []models.Metric{{Name: "avg1", Value: 2}}}}},
				{Name: "load_avg", Groups: []models.Group{{Metrics: []models.Metric{{Name: "avg1", Value: 3}}}}},
			},
		},
		{
			title: "get values of cpu_usage that is not present in storage",
			n:     3,
			name:  "cpu_usage",
		},
		{
			title: "invalid n",
			n:     -3,
		},
	}

	for _, tc := range cases {
		t.Run(tc.title, func(t *testing.T) {
			storage := NewStorage()
			storage.metrics[tc.name] = tc.metrics

			actual := storage.Get(tc.name, tc.n)
			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestSave(t *testing.T) {
	buffer = 2

	metrics := []models.Metrics{
		{Name: "load_avg", Groups: []models.Group{{Metrics: []models.Metric{{Name: "avg1", Value: 1}}}}},
		{Name: "load_avg", Groups: []models.Group{{Metrics: []models.Metric{{Name: "avg1", Value: 2}}}}},
		{Name: "load_avg", Groups: []models.Group{{Metrics: []models.Metric{{Name: "avg1", Value: 3}}}}},
		{Name: "load_avg", Groups: []models.Group{{Metrics: []models.Metric{{Name: "avg1", Value: 2}}}}},
		{Name: "load_avg", Groups: []models.Group{{Metrics: []models.Metric{{Name: "avg1", Value: 1}}}}},
	}

	storage := NewStorage()
	for i := range metrics {
		assert.NotPanics(t, func() { storage.Save(metrics[i]) })
		assert.True(t, len(storage.metrics["load_avg"]) != 0)
	}
}
