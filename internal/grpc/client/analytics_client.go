package client

import (
	"context"
	"time"

	"github.com/rainbow96bear/planet_utils/pb"
	"github.com/rainbow96bear/planet_utils/pkg/logger"
	"google.golang.org/grpc"
)

// AnalyticsClient wraps Analytics gRPC service
type AnalyticsClient struct {
	client pb.AnalyticsServiceClient
}

// NewAnalyticsClient 생성
func NewAnalyticsClient(conn *grpc.ClientConn) *AnalyticsClient {
	return &AnalyticsClient{
		client: pb.NewAnalyticsServiceClient(conn),
	}
}

// PublishEvent 단순 gRPC 호출
func (c *AnalyticsClient) PublishEvent(ctx context.Context, req *pb.PublishEventRequest) {
	if req == nil {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()

	if _, err := c.client.PublishEvent(ctx, req); err != nil {
		logger.Warnf("analytics publish failed err=%v", err)
	}
}
