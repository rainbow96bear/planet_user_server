package bootstrap

import (
	grpcclient "github.com/rainbow96bear/planet_user_server/internal/grpc/client"
	"github.com/rainbow96bear/planet_user_server/internal/repository"
	"github.com/rainbow96bear/planet_user_server/internal/resolver"
	"github.com/rainbow96bear/planet_user_server/internal/service"
	"gorm.io/gorm"
)

type Dependencies struct {
	Repos       *Repositories
	GrpcClients *grpcclient.GrpcClients
	Services    *Services
	Resolver    *resolver.Resolver
}

type Repositories struct {
	Profile *repository.ProfileRepository
}

type Services struct {
	Profile service.ProfileServiceInterface
}

func InitDependencies(db *gorm.DB) (*Dependencies, error) {
	// --- 1. Repository 초기화 ---
	profileRepo := repository.NewProfilesRepository(db)
	calendarRepo := repository.NewCalendarEventsRepository(db)
	// todoRepo := repository.NewTodosRepository(db)

	// --- 2. gRPC Clients 초기화 ---
	grpcClients, err := grpcclient.NewGrpcClients()
	if err != nil {
		return nil, err
	}

	// --- 3. Service 초기화 ---
	profileService := service.NewProfileService(db, profileRepo)
	calendarService := service.NewCalendarService(db,
		profileRepo,
		calendarRepo,
		// todoRepo,
	)

	resolver := resolver.NewResolver(profileService, calendarService)
	// DI Container 패턴
	return &Dependencies{
		Repos: &Repositories{
			Profile: profileRepo,
		},
		GrpcClients: grpcClients,
		Services: &Services{
			Profile: profileService,
		},
		Resolver: resolver,
	}, nil
}
