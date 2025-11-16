package bootstrap

import (
	"github.com/rainbow96bear/planet_user_server/internal/handler"
	"github.com/rainbow96bear/planet_user_server/internal/repository"
	"github.com/rainbow96bear/planet_user_server/internal/service"
	"github.com/rainbow96bear/planet_utils/pkg/router"
	"gorm.io/gorm"
)

func InitHandlers(db *gorm.DB) map[string]router.RouteRegistrar {
	userRepo := &repository.UsersRepository{DB: db}
	followRepo := &repository.FollowsRepository{DB: db}
	calendarRepo := &repository.CalendarRepository{DB: db}

	profileService := &service.ProfileService{
		UsersRepo: userRepo,
	}

	followService := &service.FollowService{
		UsersRepo:   userRepo,
		FollowsRepo: followRepo,
	}

	settingService := &service.SettingService{
		UsersRepo: userRepo,
	}

	calendarService := &service.CalendarService{
		CalendarRepo: calendarRepo,
		UsersRepo:    userRepo,
	}

	return map[string]router.RouteRegistrar{
		"profile":  handler.NewProfileHandler(profileService, followService),
		"user":     handler.NewUserHandler(profileService, followService),
		"setting":  handler.NewSettingHandler(settingService),
		"calendar": handler.NewCalendarHandler(calendarService),
	}
}
