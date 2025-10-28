package bootstrap

import (
	"database/sql"

	"github.com/rainbow96bear/planet_user_server/internal/handler"
	"github.com/rainbow96bear/planet_user_server/internal/repository"
	"github.com/rainbow96bear/planet_user_server/internal/service"
	"github.com/rainbow96bear/planet_utils/pkg/router"
)

func InitHandlers(db *sql.DB) map[string]router.RouteRegistrar {
	userRepo := &repository.UsersRepository{DB: db}
	followRepo := &repository.FollowsRepository{DB: db}

	profileService := &service.ProfileService{UsersRepo: userRepo}
	authService := &service.AuthService{UsersRepo: userRepo}
	followService := &service.FollowService{
		UsersRepo:   userRepo,
		FollowsRepo: followRepo,
	}

	return map[string]router.RouteRegistrar{
		"profile": handler.NewProfileHandler(profileService, authService, followService),
		"user":    handler.NewUserHandler(profileService, authService, followService),
	}
}
