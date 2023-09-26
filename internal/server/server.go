package server

import (
	"fmt"
	"net"

	"github.com/natalya-revtova/system-monitoring/internal/config"
	"github.com/natalya-revtova/system-monitoring/internal/logger"
	"github.com/natalya-revtova/system-monitoring/internal/models"
	pb "github.com/natalya-revtova/system-monitoring/pkg/api/monitoringpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Monitoring interface {
	MetricsSnapshot(m int) []models.Metrics
}

type Server struct {
	addr string
	srv  *grpc.Server
	mon  Monitoring
	pb.UnimplementedSystemMonitoringServer
	log logger.ILogger
}

func NewServer(cfg config.ServerConfig, mon Monitoring, log logger.ILogger) *Server {
	serverOptions := []grpc.ServerOption{
		grpc.Creds(insecure.NewCredentials()),
	}

	srv := grpc.NewServer(serverOptions...)

	return &Server{
		addr: fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		srv:  srv,
		mon:  mon,
		log:  log,
	}
}

func (s *Server) Start() error {
	lsn, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	pb.RegisterSystemMonitoringServer(s.srv, s)
	return s.srv.Serve(lsn)
}

func (s *Server) Stop() {
	s.srv.GracefulStop()
}
