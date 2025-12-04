package bootstrap

import (
	grpcclient "github.com/rainbow96bear/planet_user_server/internal/grpc/client"
	"gorm.io/gorm"
)

type Dependencies struct {
	Repos       *Repositories
	GrpcClients *grpcclient.GrpcClients
	Services    *Services
}

type Repositories struct {
	// User repository.UserRepository
}

type Services struct {
	// Auth  service.AuthServiceInterface
}

func InitDependencies(db *gorm.DB) (*Dependencies, error) {
	// --- 1. Repository 초기화 ---
	// userRepo := repository.NewUserRepository(db)

	// --- 2. gRPC Clients 초기화 ---
	grpcClients, err := grpcclient.NewGrpcClients()
	if err != nil {
		return nil, err
	}

	// --- 3. Service 초기화 ---
	// authService := service.NewAuthService(grpcClients)

	// DI Container 패턴
	return &Dependencies{
		// Repos: &Repositories{
		// 	User: userRepo,
		// },
		GrpcClients: grpcClients,
		Services:    &Services{
			// Auth:  authService,
		},
	}, nil
}
