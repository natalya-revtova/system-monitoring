package server

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/natalya-revtova/system-monitoring/internal/models"
	pb "github.com/natalya-revtova/system-monitoring/pkg/api/monitoringpb"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CollectMetrics(params *pb.CollectParams, stream pb.SystemMonitoring_CollectMetricsServer) error {
	clientID := newClientID()
	log := s.log.With(slog.String("client_id", clientID))

	log.Info("Client connected")

	notifyInterval := time.Duration(params.GetNotifyInterval()) * time.Second
	if notifyInterval <= 0 {
		err := errors.New("invalid NotifyInterval")
		log.Error("Validate input parameters", "error", err)
		return status.Error(codes.InvalidArgument, err.Error())
	}

	averageInterval := params.GetAverageInterval()
	if averageInterval <= 0 {
		err := errors.New("invalid AverageInterval")
		log.Error("Validate input parameters", "error", err)
		return status.Error(codes.InvalidArgument, err.Error())
	}

	ctx := stream.Context()
	ticker := time.NewTicker(notifyInterval)
	for {
		select {
		case <-ctx.Done():
			log.Info("Client disconnected")
			return nil
		case <-ticker.C:
			metrics := s.mon.MetricsSnapshot(int(averageInterval))
			if err := stream.Send(toPBModel(metrics)); err != nil {
				log.Error("Send metrics", "error", err)
			}
		}
	}
}

func toPBModel(metrics []models.Metrics) *pb.Result {
	result := make([]*pb.Metrics, 0, len(metrics))

	for i := range metrics {
		pbGroups := make([]*pb.Group, 0, len(metrics[i].Groups))
		for j := range metrics[i].Groups {
			pbLabels := make([]*pb.Label, 0, len(metrics[i].Groups[j].Labels))
			for l := range metrics[i].Groups[j].Labels {
				pbLabels = append(pbLabels, &pb.Label{
					Name:  metrics[i].Groups[j].Labels[l].Name,
					Value: metrics[i].Groups[j].Labels[l].Value,
				})
			}
			pbMetrics := make([]*pb.Metric, 0, len(metrics[i].Groups[j].Metrics))
			for m := range metrics[i].Groups[j].Metrics {
				pbMetrics = append(pbMetrics, &pb.Metric{
					Name:  metrics[i].Groups[j].Metrics[m].Name,
					Value: metrics[i].Groups[j].Metrics[m].Value,
				})
			}
			pbGroups = append(pbGroups, &pb.Group{
				Label:  pbLabels,
				Metric: pbMetrics,
			})
		}
		result = append(result, &pb.Metrics{
			Name:   metrics[i].Name,
			Groups: pbGroups,
		})
	}
	return &pb.Result{
		Metrics: result,
	}
}

func newClientID() string {
	return uuid.New().String()
}
