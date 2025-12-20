package resolver

import "github.com/rainbow96bear/planet_user_server/internal/service"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require
// here.

type Resolver struct {
	ProfileService  service.ProfileServiceInterface
	CalendarService service.CalendarServiceInterface
	TodoService     service.TodoServiceInterface
}

func NewResolver(
	profileSvc service.ProfileServiceInterface,
	calendarSvc service.CalendarServiceInterface,
	todoSvc service.TodoServiceInterface,
) *Resolver {
	return &Resolver{
		ProfileService:  profileSvc,
		CalendarService: calendarSvc,
		TodoService:     todoSvc,
	}
}
