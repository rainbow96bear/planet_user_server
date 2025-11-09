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
	// 추후 이미지 업로드 서비스 추가 시
	// ImageUploader dto.ImageUploader
}

// GetUserCalendar 사용자의 모든 캘린더 조회
func (s *CalendarService) GetUserCalendar(ctx context.Context, userUuid string) ([]*dto.CalendarInfo, error) {
	calendars, err := s.CalendarRepo.GetCalendarsByUserUuid(ctx, userUuid)
	if err != nil {
		return nil, err
	}
	return calendars, nil
}

// GetCalendarByID 특정 캘린더 조회 (권한 확인 포함)
func (s *CalendarService) GetCalendarByID(ctx context.Context, userUuid string, eventId int64) (*dto.CalendarInfo, error) {
	// 권한 확인
	isOwner, err := s.CalendarRepo.IsOwnerOfCalendar(ctx, eventId, userUuid)
	if err != nil {
		return nil, err
	}
	if !isOwner {
		return nil, fmt.Errorf("unauthorized to access this event")
	}

	// 캘린더 조회
	calendar, err := s.CalendarRepo.GetCalendarByID(ctx, eventId)
	if err != nil {
		return nil, err
	}

	return calendar, nil
}

// CreateCalendar 새 캘린더 생성
func (s *CalendarService) CreateCalendar(ctx context.Context, calendar *dto.CalendarInfo) error {
	if calendar == nil {
		return fmt.Errorf("invalid calendar data")
	}

	// 필수 필드 검증
	if calendar.Title == "" {
		return fmt.Errorf("title is required")
	}
	if calendar.StartAt == "" || calendar.EndAt == "" {
		return fmt.Errorf("start date and end date are required")
	}
	if calendar.Visibility == "" {
		return fmt.Errorf("visibility is required")
	}

	// Visibility 값 검증
	if calendar.Visibility != "public" && calendar.Visibility != "friends" && calendar.Visibility != "private" {
		return fmt.Errorf("invalid visibility value")
	}

	return s.CalendarRepo.CreateCalendar(ctx, calendar)
}

// UpdateCalendar 캘린더 수정
func (s *CalendarService) UpdateCalendar(ctx context.Context, calendar *dto.CalendarInfo) error {
	if calendar == nil {
		return fmt.Errorf("invalid data")
	}

	// 권한 확인
	isOwner, err := s.CalendarRepo.IsOwnerOfCalendar(ctx, calendar.EventID, calendar.UserUUID)
	if err != nil {
		return err
	}
	if !isOwner {
		return fmt.Errorf("unauthorized to update this event")
	}

	// 필수 필드 검증
	if calendar.Title == "" {
		return fmt.Errorf("title is required")
	}
	if calendar.StartAt == "" || calendar.EndAt == "" {
		return fmt.Errorf("start date and end date are required")
	}

	// Visibility 값 검증
	if calendar.Visibility != "public" && calendar.Visibility != "friends" && calendar.Visibility != "private" {
		return fmt.Errorf("invalid visibility value")
	}

	return s.CalendarRepo.UpdateCalendar(ctx, calendar)
}

// DeleteCalendar 캘린더 삭제
func (s *CalendarService) DeleteCalendar(ctx context.Context, userUuid string, eventId int64) error {
	// 권한 확인
	isOwner, err := s.CalendarRepo.IsOwnerOfCalendar(ctx, eventId, userUuid)
	if err != nil {
		return err
	}
	if !isOwner {
		return fmt.Errorf("unauthorized to delete this event")
	}

	// 추후 이미지 삭제 로직 추가
	// calendar, err := s.CalendarRepo.GetCalendarByID(ctx, eventId)
	// if err != nil {
	// 	return err
	// }
	// if calendar.ImageURL != "" {
	// 	_ = s.ImageUploader.Delete(calendar.ImageURL)
	// }

	return s.CalendarRepo.DeleteCalendar(ctx, eventId)
}

// GetPublicCalendarByNickname 닉네임으로 공개 캘린더 조회
func (s *CalendarService) GetPublicCalendarByNickname(ctx context.Context, nickname string) ([]*dto.CalendarInfo, error) {
	userUuid, err := s.UsersRepo.GetUserUuidByNickname(ctx, nickname)
	if err != nil {
		return nil, err
	}
	return s.CalendarRepo.GetPublicCalendarsByUserUuid(ctx, userUuid)
}

// GetAllPublicCalendars 모든 공개 캘린더 조회
func (s *CalendarService) GetAllPublicCalendars(ctx context.Context) ([]*dto.CalendarInfo, error) {
	return s.CalendarRepo.GetAllPublicCalendars(ctx)
}
