package client

import (
	"time"

	"github.com/rainbow96bear/planet_utils/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

func NewGrpcConn(target string) (*grpc.ClientConn, error) {
	logger.Infof("dialing grpc target=%s", target)

	conn, err := grpc.Dial(
		target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithTimeout(3*time.Second),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                30 * time.Second,
			Timeout:             5 * time.Second,
			PermitWithoutStream: true,
		}),
	)
	if err != nil {
		logger.Errorf("failed to dial grpc target=%s err=%v", target, err)
		return nil, err
	}

	logger.Infof("connected grpc target=%s", target)
	return conn, nil
}
