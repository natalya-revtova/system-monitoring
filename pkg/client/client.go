package client

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"

	pb "github.com/natalya-revtova/system-monitoring/pkg/api/monitoringpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	addr string
}

func NewClient(host, port string) *Client {
	return &Client{
		addr: net.JoinHostPort(host, port),
	}
}

func (c *Client) Start(ctx context.Context, notifyInterval, averageInterval int) error {
	conn, err := grpc.Dial(c.addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pb.NewSystemMonitoringClient(conn)
	stream, err := client.CollectMetrics(ctx, &pb.CollectParams{
		NotifyInterval:  int64(notifyInterval),
		AverageInterval: int64(averageInterval),
	})
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		metrics, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return err
		}
		fmt.Println(metrics)
	}
}
