package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rainbow96bear/planet_user_server/internal/models"
	"github.com/rainbow96bear/planet_user_server/internal/tx"
	"github.com/rainbow96bear/planet_utils/pkg/logger"
	"gorm.io/gorm"
)

type CalendarEventsRepository struct {
	db *gorm.DB
}

func NewCalendarEventsRepository(db *gorm.DB) *CalendarEventsRepository {
	if db == nil {
		panic("database connection is required")
	}
	return &CalendarEventsRepository{
		db: db,
	}
}

func (r *CalendarEventsRepository) getDB(ctx context.Context) *gorm.DB {
	// tx 패키지를 사용하여 Context에서 트랜잭션을 추출합니다.
	if tx := tx.GetTx(ctx); tx != nil {
		return tx.WithContext(ctx)
	}
	return r.db.WithContext(ctx) // 기본 DB 연결 반환
}

// -------------------------
// 캘린더 이벤트 생성 (Todos 포함)
// -------------------------
func (r *CalendarEventsRepository) CreateCalendarEvent(
	ctx context.Context,
	event *models.CalendarEvent,
) (*models.CalendarEvent, error) {

	db := r.getDB(ctx)

	logger.Debugf(
		"[CalendarRepo] create start user=%s title=%s",
		event.UserID,
		event.Title,
	)

	err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(event).Error; err != nil {
			logger.Errorf(
				"[CalendarRepo] insert failed user=%s err=%v",
				event.UserID,
				err,
			)
			return fmt.Errorf("failed to insert calendar event: %w", err)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	logger.Infof(
		"[CalendarRepo] created event id=%s user=%s",
		event.ID,
		event.UserID,
	)

	return event, nil
}

// -------------------------
// 단일 조회 (Todos 포함)
// -------------------------
func (r *CalendarEventsRepository) FindByID(
	ctx context.Context,
	eventID uuid.UUID,
) (*models.CalendarEvent, error) {
	db := r.getDB(ctx)

	var event models.CalendarEvent
	if err := db.
		Preload("Todos").
		First(&event, "id = ?", eventID).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to find calendar event: %w", err)
	}

	return &event, nil
}

// // -------------------------
// // 캘린더 이벤트 삭제 (Todos 포함)
// // -------------------------
func (r *CalendarEventsRepository) DeleteCalendarEvent(ctx context.Context, eventID uuid.UUID) error {
	db := r.getDB(ctx)
	logger.Infof("Deleting calendar event: %s", eventID)

	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Todos 먼저 삭제 (Foreign Key 제약 조건)
		if err := tx.Where("calendar_event_id = ?", eventID).Delete(&models.Todo{}).Error; err != nil {
			return fmt.Errorf("failed to delete todos: %w", err)
		}
		// Event 삭제
		if err := tx.Where("id = ?", eventID).Delete(&models.CalendarEvent{}).Error; err != nil {
			return fmt.Errorf("failed to delete calendar event: %w", err)
		}
		logger.Infof("Deleted calendar event %s and its todos", eventID)
		return nil
	})
}

// -------------------------
// 캘린더 이벤트 업데이트 (Todos 포함)
// -------------------------
func (r *CalendarEventsRepository) Update(
	ctx context.Context,
	event *models.CalendarEvent,
) error {

	db := r.getDB(ctx)

	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		// 🔹 CalendarEvent 업데이트
		if err := tx.Save(event).Error; err != nil {
			return err
		}

		// 🔹 기존 Todos 전체 삭제
		if err := tx.
			Where("calendar_event_id = ?", event.ID).
			Delete(&models.Todo{}).
			Error; err != nil {
			return err
		}

		// 🔹 새 Todos 삽입
		if len(event.Todos) > 0 {
			if err := tx.Create(&event.Todos).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// // ------------------------------------------
// // 조회 함수 1: 월별 뷰 (Event만, 캐시 지원)
// // ------------------------------------------

// // FindEventsWithoutTodosByVisibility: 특정 기간 동안의 Event를 Todo 없이 조회합니다.
// // CalendarService의 GetEventsWithoutTodos에서 사용됩니다. (캐싱 목적)
func (r *CalendarEventsRepository) FindEventsWithoutTodosByVisibility(
	ctx context.Context,
	UserID uuid.UUID,
	visibilities []string,
	startAt, endAt time.Time,
) ([]*models.CalendarEvent, error) {
	db := r.getDB(ctx)
	logger.Infof("Fetching events (without todos) for user=%s with visibilities=%v", UserID, visibilities)

	if len(visibilities) == 0 {
		return []*models.CalendarEvent{}, nil
	}

	var events []*models.CalendarEvent
	// 💡 Preload("Todos")를 제거하여 Todo 조인을 막습니다.
	if err := db.WithContext(ctx).
		Where("user_id = ? AND visibility IN ? AND start_at < ? AND end_at >= ?", UserID, visibilities, endAt, startAt).
		Order("start_at ASC").
		Find(&events).Error; err != nil {
		return nil, fmt.Errorf("failed to query events without todos by visibility: %w", err)
	}

	logger.Infof("Found %d calendar events (without todos) for user %s with visibility filter", len(events), UserID)
	return events, nil
}

// // ------------------------------------------
// // 조회 함수 2: 일별 뷰 (Event + Todo, 캐시 미지원)
// // ------------------------------------------

func (r *CalendarEventsRepository) FindCalendarsWithTodos(
	ctx context.Context,
	UserID uuid.UUID,
	visibilities []string,
	startAt, endAt time.Time,
) ([]*models.CalendarEvent, error) {
	db := r.getDB(ctx)
	logger.Infof("Fetching calendars (with todos) for user=%s with visibilities=%v", UserID, visibilities)

	if len(visibilities) == 0 {
		return []*models.CalendarEvent{}, nil
	}

	var events []*models.CalendarEvent
	// 💡 Preload("Todos")를 포함하여 Todo를 함께 조회합니다.
	if err := db.WithContext(ctx).
		Where("user_id = ? AND visibility IN ? AND start_at < ? AND end_at >= ?", UserID, visibilities, endAt, startAt).
		Order("start_at ASC").
		Preload("Todos").
		Find(&events).Error; err != nil {
		return nil, fmt.Errorf("failed to query calendars with todos by visibility: %w", err)
	}
	logger.Infof("Found %d calendar events (with todos) for user %s with visibility filter", len(events), UserID)
	return events, nil
}

func (r *CalendarEventsRepository) GetEventWithTodosByID(
	ctx context.Context,
	eventID uuid.UUID,
) (*models.CalendarEvent, error) {

	db := r.getDB(ctx)

	var event models.CalendarEvent
	if err := db.WithContext(ctx).
		Preload("Todos").
		First(&event, "id = ?", eventID).Error; err != nil {

		return nil, fmt.Errorf("failed to query event by ID: %w", err)
	}

	return &event, nil
}

// // -------------------------
// // 범위 조회 (visibility 없이, 전체) - 기존 함수 수정 및 유지 (Todos 포함)
// // -------------------------

// func (r *CalendarEventsRepository) FindCalendarsByUser(
// 	ctx context.Context,
// 	UserID uuid.UUID,
// 	startAt, endAt time.Time,
// ) ([]*models.CalendarEvents, error) {
// 	logger.Infof("Fetching ALL calendars for user: %s", UserID)

// 	var events []*models.CalendarEvents
// 	// Todos를 포함하여 조회합니다.
// 	if err := r.DB.WithContext(ctx).
// 		Where("user_id = ? AND start_at < ? AND end_at >= ?", UserID, endAt, startAt).
// 		Order("start_at ASC").
// 		Preload("Todos").
// 		Find(&events).Error; err != nil {
// 		return nil, fmt.Errorf("failed to query calendars: %w", err)
// 	}

// 	logger.Infof("Found %d ALL calendar events for user %s", len(events), UserID)
// 	return events, nil
// }
