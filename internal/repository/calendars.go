package repository

import (
	"context"
	"database/sql"

	"github.com/rainbow96bear/planet_user_server/dto"
)

type CalendarRepository struct {
	DB *sql.DB
}

// 내 일정 목록
func (r *CalendarRepository) GetCalendarsByUserUuid(ctx context.Context, userUuid string) ([]*dto.CalendarInfo, error) {
	query := `
		SELECT id, user_id, title, description, emoji, start_at, end_at, visibility, image_url
		FROM calendar_events
		WHERE user_id = ?
		ORDER BY start_at DESC
	`
	rows, err := r.DB.QueryContext(ctx, query, userUuid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var calendars []*dto.CalendarInfo
	for rows.Next() {
		c := &dto.CalendarInfo{}
		if err := rows.Scan(&c.EventID, &c.UserUUID, &c.Title, &c.Description, &c.Emoji, &c.StartAt, &c.EndAt, &c.Visibility, &c.ImageURL); err != nil {
			return nil, err
		}
		calendars = append(calendars, c)
	}
	return calendars, nil
}

// 일정 생성
func (r *CalendarRepository) CreateCalendar(ctx context.Context, calendar *dto.CalendarInfo) error {
	query := `
		INSERT INTO calendar_events (user_id, title, description, emoji, start_at, end_at, visibility, image_url)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.DB.ExecContext(ctx, query,
		calendar.UserUUID,
		calendar.Title,
		calendar.Description,
		calendar.Emoji,
		calendar.StartAt,
		calendar.EndAt,
		calendar.Visibility,
		calendar.ImageURL,
	)
	return err
}

// 일정 수정
func (r *CalendarRepository) UpdateCalendar(ctx context.Context, calendar *dto.CalendarInfo) error {
	query := `
		UPDATE calendar_events
		SET title=?, description=?, emoji=?, start_at=?, end_at=?, visibility=?, image_url=?
		WHERE id=? AND user_id=?
	`
	_, err := r.DB.ExecContext(ctx, query,
		calendar.Title,
		calendar.Description,
		calendar.Emoji,
		calendar.StartAt,
		calendar.EndAt,
		calendar.Visibility,
		calendar.ImageURL,
		calendar.EventID,
		calendar.UserUUID,
	)
	return err
}

// 일정 삭제
func (r *CalendarRepository) DeleteCalendar(ctx context.Context, eventId int64) error {
	query := `DELETE FROM calendar_events WHERE id = ?`
	_, err := r.DB.ExecContext(ctx, query, eventId)
	return err
}

// 소유자 확인
func (r *CalendarRepository) IsOwnerOfCalendar(ctx context.Context, eventId int64, userUuid string) (bool, error) {
	query := `SELECT COUNT(*) FROM calendar_events WHERE id = ? AND user_id = ?`
	var count int
	err := r.DB.QueryRowContext(ctx, query, eventId, userUuid).Scan(&count)
	return count > 0, err
}

// 특정 유저의 공개 일정
func (r *CalendarRepository) GetPublicCalendarsByUserUuid(ctx context.Context, userUuid string) ([]*dto.CalendarInfo, error) {
	query := `
		SELECT id, user_id, title, description, emoji, start_at, end_at, visibility, image_url
		FROM calendar_events
		WHERE user_id = ? AND visibility = 'public'
		ORDER BY start_at DESC
	`
	rows, err := r.DB.QueryContext(ctx, query, userUuid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var calendars []*dto.CalendarInfo
	for rows.Next() {
		c := &dto.CalendarInfo{}
		if err := rows.Scan(&c.EventID, &c.UserUUID, &c.Title, &c.Description, &c.Emoji, &c.StartAt, &c.EndAt, &c.Visibility, &c.ImageURL); err != nil {
			return nil, err
		}
		calendars = append(calendars, c)
	}
	return calendars, nil
}

// 전체 공개 일정
func (r *CalendarRepository) GetAllPublicCalendars(ctx context.Context) ([]*dto.CalendarInfo, error) {
	query := `
		SELECT id, user_id, title, description, emoji, start_at, end_at, visibility, image_url
		FROM calendar_events
		WHERE visibility = 'public'
		ORDER BY start_at DESC
	`
	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var calendars []*dto.CalendarInfo
	for rows.Next() {
		c := &dto.CalendarInfo{}
		if err := rows.Scan(&c.EventID, &c.UserUUID, &c.Title, &c.Description, &c.Emoji, &c.StartAt, &c.EndAt, &c.Visibility, &c.ImageURL); err != nil {
			return nil, err
		}
		calendars = append(calendars, c)
	}
	return calendars, nil
}
