package bootstrap

import (
	"github.com/rainbow96bear/planet_user_server/internal/handler"
	"github.com/rainbow96bear/planet_user_server/internal/repository"
	"github.com/rainbow96bear/planet_user_server/internal/service"
	"github.com/rainbow96bear/planet_utils/pkg/router"
	"gorm.io/gorm"
)

func InitHandlers(db *gorm.DB) map[string]router.RouteRegistrar {
	profilesRepo := &repository.ProfilesRepository{DB: db}
	followRepo := &repository.FollowsRepository{DB: db}
	calendarEventsRepo := &repository.CalendarEventsRepository{DB: db}

	profileService := &service.ProfileService{
		ProfilesRepo: profilesRepo,
	}

	followService := &service.FollowService{
		ProfilesRepo: profilesRepo,
		FollowsRepo:  followRepo,
	}

	calendarService := &service.CalendarService{
		CalendarEventsRepo: calendarEventsRepo,
	}

	return map[string]router.RouteRegistrar{
		"profile":  handler.NewProfileHandler(profileService, followService),
		"follow":   handler.NewFollowHandler(profileService, followService),
		"calendar": handler.NewCalendarHandler(calendarService),
	}
}
