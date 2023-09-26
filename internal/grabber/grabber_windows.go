//go:build windows
// +build windows

package grabber

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	"github.com/natalya-revtova/system-monitoring/internal/config"
	"github.com/natalya-revtova/system-monitoring/internal/logger"
	"github.com/natalya-revtova/system-monitoring/internal/models"
)

type WindowsGrabber struct {
	options []string
	log     logger.ILogger
}

func NewGrabber(options []string, log logger.ILogger) *WindowsGrabber {
	return &WindowsGrabber{
		options: options,
		log:     log,
	}
}

func GetOptions(cfg config.MetricsConfig) []string {
	options := make([]string, 0)
	if cfg.CPUUsg {
		options = append(options, models.CPUStatOption)
	}
	return options
}

func (g *WindowsGrabber) Grab(results chan models.Metrics) {
	wg := &sync.WaitGroup{}
	wg.Add(len(g.options))

	for _, option := range g.options {
		go func(option string) {
			defer wg.Done()

			switch option {
			case models.CPUStatOption:
				g.cpuStat(results)
			}
		}(option)
	}

	wg.Wait()
}

func (g *WindowsGrabber) cpuStat(results chan<- models.Metrics) {
	output, err := exec.Command("top",
		`Processor Information(_Total)\% Privileged Time`,
		`Processor Information(_Total)\% User Time`,
		`Processor Information(_Total)\% Idle Time`,
		"-sc", "1").Output()
	if err != nil {
		g.log.Error("Collect CPU usage", "error", err)
		return
	}
	fields := strings.Fields(string(output))

	results <- models.Metrics{
		Name: models.CPUStatOption,
		Groups: []models.Group{
			{
				Metrics: []models.Metric{
					{Name: "user_mode", Value: g.parseValue(fields[2])},
					{Name: "system_mode", Value: g.parseValue(fields[1])},
					{Name: "idle", Value: g.parseValue(fields[3])},
				},
			},
		},
	}
}

func (g *WindowsGrabber) parseValue(value string) float64 {
	parsed, err := strconv.ParseFloat(strings.Trim(value, "\""), 64)
	if err != nil {
		g.log.Error("Parse string value", "error", err)
		return 0.0
	}
	return parsed
}
