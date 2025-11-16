package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/rainbow96bear/planet_utils/model"
	"github.com/rainbow96bear/planet_utils/pkg/logger"
	"gorm.io/gorm"
)

type CalendarRepository struct {
	DB *gorm.DB
}

// 특정 유저의 일정 중 공개범위(visibilities)에 해당하고,
// startDate ~ endDate 범위 안에 속하는 일정 조회
func (r *CalendarRepository) FindCalendarsByVisibility(
	ctx context.Context,
	userUUID string,
	visibilities []string,
	startAt, endAt time.Time,
) ([]*model.Calendar, error) {
	logger.Infof("start to find calendar events uuid : %s", userUUID)
	defer logger.Infof("end to find calendar events uuid : %s", userUUID)

	if len(visibilities) == 0 {
		return nil, nil
	}

	var results []*model.Calendar
	if err := r.DB.WithContext(ctx).
		Where("user_uuid = ? AND visibility IN ? AND start_at < ? AND end_at >= ?", userUUID, visibilities, endAt, startAt).
		Order("start_at ASC").
		Preload("Todos").
		Find(&results).Error; err != nil {
		return nil, fmt.Errorf("query calendars: %w", err)
	}

	for i, cal := range results {
		logger.Debugf("Calendar %d: %+v", i, *cal)            // 포인터를 역참조해서 내용 출력
		logger.Debugf("Calendar %d Todos: %+v", i, cal.Todos) // Todos 내용도 출력
	}
	return results, nil
}

// Calendar + Todos 생성
func (r *CalendarRepository) CreateCalendarWithTodos(ctx context.Context, cal *model.Calendar) error {
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(cal).Error; err != nil {
			return fmt.Errorf("insert calendar: %w", err)
		}

		for i := range cal.Todos {
			cal.Todos[i].EventID = cal.EventID
		}
		if len(cal.Todos) > 0 {
			if err := tx.Create(&cal.Todos).Error; err != nil {
				return fmt.Errorf("insert todos: %w", err)
			}
		}

		return nil
	})
}

// ID로 단건 조회
func (r *CalendarRepository) FindByID(ctx context.Context, eventID uint64) (*model.Calendar, error) {
	var cal model.Calendar
	if err := r.DB.WithContext(ctx).
		Preload("Todos").
		First(&cal, "event_id = ?", eventID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("find calendar by id: %w", err)
	}
	return &cal, nil
}

// Calendar 삭제
func (r *CalendarRepository) DeleteCalendar(ctx context.Context, eventID uint64) error {
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("event_id = ?", eventID).Delete(&model.Todo{}).Error; err != nil {
			return fmt.Errorf("delete todos: %w", err)
		}
		if err := tx.Where("event_id = ?", eventID).Delete(&model.Calendar{}).Error; err != nil {
			return fmt.Errorf("delete calendar: %w", err)
		}
		return nil
	})
}
