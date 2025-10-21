package grpc_client

import (
	"context"
	"time"

	"github.com/rainbow96bear/planet_db_server/logger"
	pb "github.com/rainbow96bear/planet_proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type DBClient struct {
	conn   *grpc.ClientConn
	client pb.UserServiceClient
}

func NewDBClient(addr string) (*DBClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Errorf("grpc server did not connect: %v", err)
		return nil, err
	}
	return &DBClient{
		conn:   conn,
		client: pb.NewUserServiceClient(conn),
	}, nil
}

func (d *DBClient) ReqGetUserInfoByNickname(userInfo *pb.UserInfo) (*pb.UserInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return d.client.GetUserInfoByNickname(ctx, userInfo)
}

func (d *DBClient) ReqUpdateUserProfile(userInfo *pb.UserInfo) (*pb.UserInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return d.client.GetUserInfoByNickname(ctx, userInfo)
}
