//go:build linux
// +build linux

package grabber

import (
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	"github.com/natalya-revtova/system-monitoring/internal/config"
	"github.com/natalya-revtova/system-monitoring/internal/logger"
	"github.com/natalya-revtova/system-monitoring/internal/models"
)

type LinuxGrabber struct {
	options []string
	log     logger.ILogger
}

func NewGrabber(options []string, log logger.ILogger) LinuxGrabber {
	return LinuxGrabber{
		options: options,
		log:     log,
	}
}

func GetOptions(cfg config.MetricsConfig) []string {
	options := make([]string, 0)
	if cfg.LoadAvg {
		options = append(options, models.LoadAverageOption)
	}
	if cfg.CPUUsg {
		options = append(options, models.CPUStatOption)
	}
	if cfg.DiskInfo {
		options = append(options, models.DiskStatOption)
	}
	return options
}

func (g LinuxGrabber) Grab(results chan models.Metrics) {
	wg := &sync.WaitGroup{}
	wg.Add(len(g.options))

	for _, option := range g.options {
		go func(option string) {
			defer wg.Done()

			switch option {
			case models.LoadAverageOption:
				g.loadAverage(results)
			case models.CPUStatOption:
				g.cpuStat(results)
			case models.DiskStatOption:
				g.diskStat(results)
			}
		}(option)
	}

	wg.Wait()
}

func (g LinuxGrabber) loadAverage(results chan<- models.Metrics) {
	data, err := os.ReadFile("/proc/loadavg")
	if err != nil {
		g.log.Error("Read load average file", "error", err)
		return
	}
	fields := strings.Fields(string(data))

	results <- models.Metrics{
		Name: models.LoadAverageOption,
		Groups: []models.Group{
			{
				Metrics: []models.Metric{
					{Name: "1_min", Value: g.parseValue(fields[0])},
					{Name: "5_min", Value: g.parseValue(fields[1])},
					{Name: "15_min", Value: g.parseValue(fields[2])},
				},
			},
		},
	}
}

func (g LinuxGrabber) cpuStat(results chan<- models.Metrics) {
	output, err := exec.Command("top", "-b", "-n1").Output()
	if err != nil {
		g.log.Error("Collect CPU usage", "error", err)
		return
	}

	var cpuInfo string
	lines := strings.Split(string(output), "\n")
	for i := range lines {
		if strings.Contains(lines[i], "Cpu") {
			cpuInfo = lines[i]
			break
		}
	}
	fields := strings.Fields(cpuInfo)

	results <- models.Metrics{
		Name: models.CPUStatOption,
		Groups: []models.Group{
			{
				Metrics: []models.Metric{
					{Name: "user_mode", Value: g.parseValue(fields[1])},
					{Name: "system_mode", Value: g.parseValue(fields[3])},
					{Name: "idle", Value: g.parseValue(fields[7])},
				},
			},
		},
	}
}

func (g LinuxGrabber) diskStat(results chan<- models.Metrics) {
	const columnsInDFOutput = 7

	dfOutput := g.getDiskInfo()
	dfOutputInode := g.getDiskInodeInfo()

	groups := make([]models.Group, 0, len(dfOutput))
	for i := range dfOutput {
		diskFields := strings.Fields(dfOutput[i])
		inodeFields := strings.Fields(dfOutputInode[i])

		if len(diskFields) < columnsInDFOutput || len(inodeFields) < columnsInDFOutput {
			continue
		}

		groups = append(groups, models.Group{
			Labels: []models.Label{
				{Name: "filesystem", Value: diskFields[0]},
				{Name: "type", Value: diskFields[1]},
				{Name: "mounted_on", Value: diskFields[6]},
			},
			Metrics: []models.Metric{
				{Name: "disk_used", Value: g.parseValue(diskFields[3]) / (1024 * 1024)}, // convert to MB
				{Name: "disk_usage", Value: g.calculateUsage(diskFields[2], diskFields[3])},
				{Name: "inode_used", Value: g.parseValue(inodeFields[3])},
				{Name: "inode_usage", Value: g.calculateUsage(inodeFields[2], inodeFields[3])},
			},
		})
	}

	results <- models.Metrics{
		Name:   models.DiskStatOption,
		Groups: groups,
	}
}

func (g LinuxGrabber) calculateUsage(total, used string) float64 {
	if g.parseValue(total) == 0.0 {
		return 0.0
	}
	return g.parseValue(used) * 100 / g.parseValue(total)
}

func (g LinuxGrabber) getDiskInfo() []string {
	dfCmd := exec.Command("df", "-T", "-B1", "--exclude-type=tmpfs")
	res, err := dfCmd.Output()
	if err != nil {
		g.log.Error("Collect disk info", "error", err)
		return nil
	}
	return strings.Split(string(res), "\n")[1:]
}

func (g LinuxGrabber) getDiskInodeInfo() []string {
	dfCmd := exec.Command("df", "-T", "-i", "--exclude-type=tmpfs")
	res, err := dfCmd.Output()
	if err != nil {
		g.log.Error("Collect inodes info", "error", err)
		return nil
	}
	return strings.Split(string(res), "\n")[1:]
}

func (g LinuxGrabber) parseValue(value string) float64 {
	parsed, err := strconv.ParseFloat(floatFormat(value), 64)
	if err != nil {
		g.log.Error("Parse string value", "error", err)
		return 0.0
	}
	return parsed
}

func floatFormat(in string) string {
	return strings.Join(strings.Split(in, ","), ".")
}
