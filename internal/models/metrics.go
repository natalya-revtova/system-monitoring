package models

const (
	LoadAverageOption = "load_average"
	CPUStatOption     = "cpu_usage"
	DiskStatOption    = "disk_usage"
)

type (
	Metric struct {
		Name  string
		Value float64
	}

	Label struct {
		Name  string
		Value string
	}

	Group struct {
		Labels  []Label
		Metrics []Metric
	}

	Metrics struct {
		Name   string
		Groups []Group
	}
)
