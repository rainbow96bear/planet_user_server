package resolver

import "github.com/rainbow96bear/planet_user_server/internal/service"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require
// here.

type Resolver struct {
	ProfileService  *service.ProfileService
	CalendarService *service.CalendarService
	TodoService     *service.TodoService
}

func NewResolver(
	profileSvc *service.ProfileService,
	calendarSvc *service.CalendarService,
	todoSvc *service.TodoService,
) *Resolver {
	return &Resolver{
		ProfileService:  profileSvc,
		CalendarService: calendarSvc,
		TodoService:     todoSvc,
	}
}
