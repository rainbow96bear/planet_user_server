package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rainbow96bear/planet_utils/models"
	"github.com/rainbow96bear/planet_utils/pkg/logger"
	"gorm.io/gorm"
)

type CalendarEventsRepository struct {
	DB *gorm.DB
}

// -------------------------
// 트랜잭션 시작
// -------------------------
func (r *CalendarEventsRepository) BeginTx(ctx context.Context) (*gorm.DB, error) {
	logger.Infof("starting transaction for CalendarEventsRepository")
	tx := r.DB.WithContext(ctx).Begin()
	if tx.Error != nil {
		logger.Errorf("failed to start transaction: %v", tx.Error)
		return nil, tx.Error
	}
	logger.Infof("transaction started successfully")
	return tx, nil
}

// -------------------------
// 캘린더 이벤트 생성 (Todos 포함)
// -------------------------
func (r *CalendarEventsRepository) CreateCalendarEvent(ctx context.Context, event *models.CalendarEvents) error {
	logger.Infof("Creating calendar event for user: %s", event.UserID)

	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Calendar Event 삽입
		if err := tx.Create(event).Error; err != nil {
			return fmt.Errorf("failed to insert calendar event: %w", err)
		}

		// 2. Todos 삽입 전, EventID와 ID 설정/재설정
		if len(event.Todos) > 0 {
			for i := range event.Todos {
				// 🌟🌟🌟 핵심 수정: 각 Todo 항목에 새로운 UUID를 강제 할당합니다. 🌟🌟🌟
				// DTO 변환 과정에서 잘못된 ID가 할당되었거나 재사용되었을 가능성 방지.
				event.Todos[i].ID = uuid.New()

				// EventID 설정 (이벤트 삽입 후 생성된 event.ID 사용)
				event.Todos[i].EventID = event.ID
			}

			// Todos 목록 삽입
			if err := tx.Create(&event.Todos).Error; err != nil {
				// 🚨 로그에서 발생한 오류 위치
				return fmt.Errorf("failed to insert todos: %w", err)
			}
		}

		logger.Infof("Successfully created calendar event %s with %d todos", event.ID, len(event.Todos))
		return nil
	})
}

// -------------------------
// 단일 조회
// -------------------------
func (r *CalendarEventsRepository) FindByID(ctx context.Context, eventID uuid.UUID) (*models.CalendarEvents, error) {
	var event models.CalendarEvents
	if err := r.DB.WithContext(ctx).
		Preload("Todos").
		First(&event, "id = ?", eventID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find calendar event: %w", err)
	}
	return &event, nil
}

// -------------------------
// 캘린더 이벤트 삭제 (Todos 포함)
// -------------------------
func (r *CalendarEventsRepository) DeleteCalendarEvent(ctx context.Context, eventID uuid.UUID) error {
	logger.Infof("Deleting calendar event: %s", eventID)

	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("event_id = ?", eventID).Delete(&models.Todos{}).Error; err != nil {
			return fmt.Errorf("failed to delete todos: %w", err)
		}
		if err := tx.Where("id = ?", eventID).Delete(&models.CalendarEvents{}).Error; err != nil {
			return fmt.Errorf("failed to delete calendar event: %w", err)
		}
		logger.Infof("Deleted calendar event %s and its todos", eventID)
		return nil
	})
}

// -------------------------
// 캘린더 이벤트 업데이트 (Todos 포함)
// -------------------------
func (r *CalendarEventsRepository) UpdateCalendarEvent(ctx context.Context, event *models.CalendarEvents) error {
	logger.Infof("[UpdateCalendar] eventID=%s", event.ID)

	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// CalendarEvent 업데이트
		if err := tx.Save(event).Error; err != nil {
			logger.Errorf("[UpdateCalendar] failed to update event: %v", err)
			return fmt.Errorf("failed to update calendar event: %w", err)
		}

		// 기존 Todos 삭제 후 새로 삽입
		if err := tx.Where("event_id = ?", event.ID).Delete(&models.Todos{}).Error; err != nil {
			logger.Errorf("[UpdateCalendar] failed to delete old todos: %v", err)
			return fmt.Errorf("failed to delete old todos: %w", err)
		}
		for i := range event.Todos {
			event.Todos[i].EventID = event.ID
		}
		if len(event.Todos) > 0 {
			if err := tx.Create(&event.Todos).Error; err != nil {
				logger.Errorf("[UpdateCalendar] failed to insert new todos: %v", err)
				return fmt.Errorf("failed to insert new todos: %w", err)
			}
		}

		logger.Infof("[UpdateCalendar] successfully updated eventID=%s with %d todos", event.ID, len(event.Todos))
		return nil
	})
}

// -------------------------
// 특정 기간 + visibility 조회
// -------------------------
func (r *CalendarEventsRepository) FindCalendarsByVisibility(
	ctx context.Context,
	UserID uuid.UUID,
	visibilities []string,
	startAt, endAt time.Time,
) ([]*models.CalendarEvents, error) {
	logger.Infof("Fetching calendars for user=%s with visibilities=%v", UserID, visibilities)

	if len(visibilities) == 0 {
		return []*models.CalendarEvents{}, nil
	}

	var events []*models.CalendarEvents
	if err := r.DB.WithContext(ctx).
		Where("user_id = ? AND visibility IN ? AND start_at < ? AND end_at >= ?", UserID, visibilities, endAt, startAt).
		Order("start_at ASC").
		Preload("Todos").
		Find(&events).Error; err != nil {
		return nil, fmt.Errorf("failed to query calendars by visibility: %w", err)
	}

	logger.Infof("Found %d calendar events for user %s with visibility filter", len(events), UserID)
	return events, nil
}

// -------------------------
// 범위 조회 (visibility 없이, 전체)
func (r *CalendarEventsRepository) FindCalendarsByUser(
	ctx context.Context,
	UserID uuid.UUID,
	startAt, endAt time.Time,
) ([]*models.CalendarEvents, error) {
	logger.Infof("Fetching calendars for user: %s", UserID)

	var events []*models.CalendarEvents
	if err := r.DB.WithContext(ctx).
		Where("user_id = ? AND start_time < ? AND end_time >= ?", UserID, endAt, startAt).
		Order("start_time ASC").
		Preload("Todos").
		Find(&events).Error; err != nil {
		return nil, fmt.Errorf("failed to query calendars: %w", err)
	}

	logger.Infof("Found %d calendar events for user %s", len(events), UserID)
	return events, nil
}
