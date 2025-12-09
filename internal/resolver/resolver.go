package resolver

import "github.com/rainbow96bear/planet_user_server/internal/service"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require
// here.

type Resolver struct {
	ProfileService service.ProfileServiceInterface
}

func NewResolver(
	profileSvc service.ProfileServiceInterface,
) *Resolver {
	return &Resolver{
		ProfileService: profileSvc,
	}
}
