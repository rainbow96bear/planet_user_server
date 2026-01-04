package bootstrap

import (
	"github.com/rainbow96bear/planet_user_server/internal/grpc/client"
	"github.com/rainbow96bear/planet_user_server/internal/repository"
	"github.com/rainbow96bear/planet_user_server/internal/resolver"
	"github.com/rainbow96bear/planet_user_server/internal/service"
	"gorm.io/gorm"
)

/*
Dependencies
👉 애플리케이션 전체 의존성 컨테이너
👉 외부에서는 "조립 결과"만 본다
*/
type Dependencies struct {
	Repos    *Repositories
	Services *Services
	Resolver *resolver.Resolver
}

// --------------------
// Repositories
// --------------------

type Repositories struct {
	Profile  *repository.ProfileRepository
	Calendar *repository.CalendarEventsRepository
	Todo     *repository.TodosRepository
	Follow   *repository.FollowsRepository
}

// --------------------
// Services
// --------------------

type Services struct {
	Profile  *service.ProfileService
	Calendar *service.CalendarService
	Todo     *service.TodoService
}

// --------------------
// InitDependencies
// --------------------

func InitDependencies(db *gorm.DB) (*Dependencies, error) {

	// 1️⃣ Repositories
	repos := initRepositories(db)

	// 2️⃣ gRPC clients (infra)
	grpcClients, err := client.NewGrpcClients()
	if err != nil {
		return nil, err
	}

	// 3️⃣ Services
	services := initServices(db, repos, grpcClients)

	// 4️⃣ Resolver (API entry point)
	res := resolver.NewResolver(
		services.Profile,
		services.Calendar,
		services.Todo,
	)

	return &Dependencies{
		Repos:    repos,
		Services: services,
		Resolver: res,
	}, nil
}

// --------------------
// internal helpers
// --------------------

func initRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		Profile:  repository.NewProfilesRepository(db),
		Calendar: repository.NewCalendarEventsRepository(db),
		Todo:     repository.NewTodosRepository(db),
		Follow:   repository.NewFollowsRepository(db),
	}
}

func initServices(
	db *gorm.DB,
	repos *Repositories,
	grpcClients *client.GrpcClients,
) *Services {

	profileSvc := service.NewProfileService(db, repos.Profile)

	calendarSvc := service.NewCalendarService(
		db,
		repos.Calendar,
		repos.Follow,
	)

	todoSvc := service.NewTodoService(
		db,
		repos.Todo,
		grpcClients.Analytics, // ✅ port만 주입
	)

	return &Services{
		Profile:  profileSvc,
		Calendar: calendarSvc,
		Todo:     todoSvc,
	}
}
