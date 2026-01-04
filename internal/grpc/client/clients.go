package client

import (
	"github.com/rainbow96bear/planet_user_server/config"
)

type GrpcClients struct {
	Analytics *AnalyticsClient
}

func NewGrpcClients() (*GrpcClients, error) {
	conn, err := NewGrpcConn(config.ANALYTICS_GRPC_SERVER_ADDR)
	if err != nil {
		return nil, err
	}

	analyticsClient := NewAnalyticsClient(conn)
	return &GrpcClients{
		Analytics: analyticsClient,
	}, nil
}
