package config

type MetricsConfig struct {
	LoadAvg  bool `toml:"load_avg"`
	CPUUsg   bool `toml:"cpu_usage"`
	DiskInfo bool `toml:"disk_info"`
}
