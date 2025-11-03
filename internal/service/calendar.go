package service

import (
	"context"
	"fmt"

	"github.com/rainbow96bear/planet_user_server/dto"
	"github.com/rainbow96bear/planet_user_server/internal/repository"
)

type CalendarService struct {
	CalendarRepo *repository.CalendarRepository
	UsersRepo    *repository.UsersRepository
}

func (s *CalendarService) GetUserCalendar(ctx context.Context, userUuid string) ([]*dto.CalendarInfo, error) {
	calendars, err := s.CalendarRepo.GetCalendarsByUserUuid(ctx, userUuid)
	if err != nil {
		return nil, err
	}
	return calendars, nil
}

func (s *CalendarService) CreateCalendar(ctx context.Context, calendar *dto.CalendarInfo) error {
	if calendar == nil {
		return fmt.Errorf("invalid calendar data")
	}
	return s.CalendarRepo.CreateCalendar(ctx, calendar)
}

func (s *CalendarService) UpdateCalendar(ctx context.Context, calendar *dto.CalendarInfo) error {
	if calendar == nil {
		return fmt.Errorf("invalid data")
	}

	isOwner, err := s.CalendarRepo.IsOwnerOfCalendar(ctx, calendar.EventID, calendar.UserUUID)
	if err != nil {
		return err
	}
	if !isOwner {
		return fmt.Errorf("unauthorized to update this event")
	}
	return s.CalendarRepo.UpdateCalendar(ctx, calendar)
}

func (s *CalendarService) DeleteCalendar(ctx context.Context, userUuid string, eventId int64) error {
	isOwner, err := s.CalendarRepo.IsOwnerOfCalendar(ctx, eventId, userUuid)
	if err != nil {
		return err
	}
	if !isOwner {
		return fmt.Errorf("unauthorized to delete this event")
	}
	return s.CalendarRepo.DeleteCalendar(ctx, eventId)
}

func (s *CalendarService) GetPublicCalendarByNickname(ctx context.Context, nickname string) ([]*dto.CalendarInfo, error) {
	userUuid, err := s.UsersRepo.GetUserUuidByNickname(ctx, nickname)
	if err != nil {
		return nil, err
	}
	return s.CalendarRepo.GetPublicCalendarsByUserUuid(ctx, userUuid)
}

func (s *CalendarService) GetAllPublicCalendars(ctx context.Context) ([]*dto.CalendarInfo, error) {
	return s.CalendarRepo.GetAllPublicCalendars(ctx)
}
