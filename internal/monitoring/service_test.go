package monitoring

import (
	"context"
	"testing"
	"time"

	"github.com/natalya-revtova/system-monitoring/internal/logger"
	"github.com/natalya-revtova/system-monitoring/internal/models"
	"github.com/natalya-revtova/system-monitoring/internal/monitoring/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

var (
	metrics = []models.Metrics{
		{
			Name: "disk_usage",
			Groups: []models.Group{
				{
					Metrics: []models.Metric{
						{Name: "disk_used", Value: 1},
						{Name: "disk_usage", Value: 2},
						{Name: "inode_used", Value: 3},
						{Name: "inode_usage", Value: 4},
					},
					Labels: []models.Label{
						{Name: "filesystem", Value: "/dev/mapper/vgubuntu-root"},
						{Name: "type", Value: "ext4"},
						{Name: "mounted_on", Value: "/"},
					},
				},
				{
					Metrics: []models.Metric{
						{Name: "disk_used", Value: 5},
						{Name: "disk_usage", Value: 6},
						{Name: "inode_used", Value: 7},
						{Name: "inode_usage", Value: 8},
					},
					Labels: []models.Label{
						{Name: "filesystem", Value: "/dev/nvme0n1p2"},
						{Name: "type", Value: "ext4"},
						{Name: "mounted_on", Value: "/boot"},
					},
				},
			},
		},
		{
			Name: "disk_usage",
			Groups: []models.Group{
				{
					Metrics: []models.Metric{
						{Name: "disk_used", Value: 9},
						{Name: "disk_usage", Value: 10},
						{Name: "inode_used", Value: 11},
						{Name: "inode_usage", Value: 12},
					},
					Labels: []models.Label{
						{Name: "filesystem", Value: "/dev/mapper/vgubuntu-root"},
						{Name: "type", Value: "ext4"},
						{Name: "mounted_on", Value: "/"},
					},
				},
				{
					Metrics: []models.Metric{
						{Name: "disk_used", Value: 13},
						{Name: "disk_usage", Value: 14},
						{Name: "inode_used", Value: 15},
						{Name: "inode_usage", Value: 16},
					},
					Labels: []models.Label{
						{Name: "filesystem", Value: "/dev/nvme0n1p2"},
						{Name: "type", Value: "ext4"},
						{Name: "mounted_on", Value: "/boot"},
					},
				},
			},
		},
	}

	summary = []models.Metrics{
		{
			Name: "disk_usage",
			Groups: []models.Group{
				{
					Metrics: []models.Metric{
						{Name: "disk_used", Value: 5},
						{Name: "disk_usage", Value: 6},
						{Name: "inode_used", Value: 7},
						{Name: "inode_usage", Value: 8},
					},
					Labels: []models.Label{
						{Name: "filesystem", Value: "/dev/mapper/vgubuntu-root"},
						{Name: "type", Value: "ext4"},
						{Name: "mounted_on", Value: "/"},
					},
				},
				{
					Metrics: []models.Metric{
						{Name: "disk_used", Value: 9},
						{Name: "disk_usage", Value: 10},
						{Name: "inode_used", Value: 11},
						{Name: "inode_usage", Value: 12},
					},
					Labels: []models.Label{
						{Name: "filesystem", Value: "/dev/nvme0n1p2"},
						{Name: "type", Value: "ext4"},
						{Name: "mounted_on", Value: "/boot"},
					},
				},
			},
		},
	}
)

func TestService(t *testing.T) {
	defer goleak.VerifyNone(t)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	grabberMock := mocks.NewGrabber(t)
	storageMock := mocks.NewStorage(t)

	grabberMock.On("Grab", mock.Anything).Return().Times(3)

	_ = NewService(ctx, grabberMock, storageMock, []string{models.DiskStatOption}, logger.NewMock())
	<-ctx.Done()
}

func TestMetricsCollection(t *testing.T) {
	defer goleak.VerifyNone(t)

	grabberMock := mocks.NewGrabber(t)
	storageMock := mocks.NewStorage(t)

	svc := Service{
		grabber: grabberMock,
		storage: storageMock,
	}

	results := make(chan models.Metrics, 1)
	results <- metrics[0]

	grabberMock.On("Grab", results).Return().Once()
	storageMock.On("Save", metrics[0]).Return().Once()

	svc.collect(results)
}

func TestAverageCalculation(t *testing.T) {
	grabberMock := mocks.NewGrabber(t)
	storageMock := mocks.NewStorage(t)

	svc := Service{
		grabber: grabberMock,
		storage: storageMock,
	}

	svc.options = []string{models.DiskStatOption}
	storageMock.On("Get", models.DiskStatOption, 2).Return(nil, true).Once()
	storageMock.On("Get", models.DiskStatOption, 2).Return(metrics, true).Once()

	actual := svc.MetricsSnapshot(2)
	require.Equal(t, summary, actual)
}
