package resolver

import grpcclient "github.com/rainbow96bear/planet_user_server/internal/grpc/client"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require
// here.

type Resolver struct {
	GrpcClients *grpcclient.GrpcClients
	// AuthService *auth.Service
}

func NewResolver(gc *grpcclient.GrpcClients) *Resolver {
	return &Resolver{
		GrpcClients: gc,
	}
}
