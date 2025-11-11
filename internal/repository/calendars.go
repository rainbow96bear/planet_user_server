package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/rainbow96bear/planet_utils/model"
)

type CalendarRepository struct {
	DB *sql.DB
}

// FindCalendarsByUserAndDateRange
// 특정 유저의 일정 중 공개범위(visibilities)에 해당하고,
// startDate ~ endDate 범위 안에 속하는 일정 조회
func (r *CalendarRepository) FindCalendarsByVisibility(
	ctx context.Context,
	userUUID string,
	visibilities []string,
	startDate, endDate time.Time,
) ([]*model.Calendar, error) {
	if len(visibilities) == 0 {
		return nil, nil
	}

	placeholders := strings.Repeat("?,", len(visibilities))
	placeholders = placeholders[:len(placeholders)-1]

	query := fmt.Sprintf(`
		SELECT event_id, user_uuid, title, start_at, end_at, emoji, visibility, created_at, updated_at
		FROM calendars
		WHERE user_uuid = ?
		  AND visibility IN (%s)
		  AND start_at < ? AND end_at >= ?
		  AND status IN ('active','completed')
		ORDER BY start_at ASC
	`, placeholders)

	args := []any{userUUID}
	for _, v := range visibilities {
		args = append(args, v)
	}
	args = append(args, endDate, startDate)

	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query calendars: %w", err)
	}
	defer rows.Close()

	var results []*model.Calendar
	for rows.Next() {
		var c model.Calendar
		if err := rows.Scan(
			&c.EventID, &c.UserUUID, &c.Title, &c.StartAt, &c.EndAt,
			&c.Emoji, &c.Visibility, &c.CreatedAt, &c.UpdatedAt,
		); err != nil {
			return nil, err
		}
		results = append(results, &c)
	}

	return results, nil
}

func (r *CalendarRepository) CreateCalendarWithTodos(ctx context.Context, cal *model.Calendar) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// 1️⃣ Calendar 생성
	queryCal := `
		INSERT INTO calendars (user_uuid, title, description, emoji, start_at, end_at, visibility, image_url, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
		RETURNING event_id, created_at, updated_at
	`
	row := tx.QueryRowContext(ctx, queryCal,
		cal.UserUUID, cal.Title, cal.Description, cal.Emoji,
		cal.StartAt, cal.EndAt, cal.Visibility, cal.ImageURL,
	)
	if err := row.Scan(&cal.EventID, &cal.CreatedAt, &cal.UpdatedAt); err != nil {
		tx.Rollback()
		return fmt.Errorf("insert calendar: %w", err)
	}

	// 2️⃣ Todos가 있으면 EventID 연결 후 삽입
	if len(cal.Todos) > 0 {
		values := make([]string, 0, len(cal.Todos))
		args := make([]any, 0, len(cal.Todos)*3)
		for _, t := range cal.Todos {
			values = append(values, "(?, ?, ?)")
			args = append(args, cal.EventID, t.Content, t.Done)
		}
		queryTodo := fmt.Sprintf(`
			INSERT INTO todos (event_id, content, done)
			VALUES %s
		`, strings.Join(values, ","))
		if _, err := tx.ExecContext(ctx, queryTodo, args...); err != nil {
			tx.Rollback()
			return fmt.Errorf("insert todos: %w", err)
		}
	}

	// 3️⃣ Commit
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("transaction commit failed: %w", err)
	}

	return nil
}

func (r *CalendarRepository) FindByID(ctx context.Context, eventID uint64) (*model.Calendar, error) {
	query := `
		SELECT event_id, user_uuid, title, description, emoji, start_at, end_at, visibility, image_url, created_at, updated_at
		FROM calendars
		WHERE event_id = ? AND status IN ('active','completed')
	`
	row := r.DB.QueryRowContext(ctx, query, eventID)

	var cal model.Calendar
	if err := row.Scan(
		&cal.EventID, &cal.UserUUID, &cal.Title, &cal.Description,
		&cal.Emoji, &cal.StartAt, &cal.EndAt, &cal.Visibility,
		&cal.ImageURL, &cal.CreatedAt, &cal.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("find calendar by id: %w", err)
	}

	// 관련 Todos 조회
	todoQuery := `SELECT id, event_id, content, done FROM todos WHERE event_id = ?`
	rows, err := r.DB.QueryContext(ctx, todoQuery, cal.EventID)
	if err != nil {
		return nil, fmt.Errorf("query todos: %w", err)
	}
	defer rows.Close()

	var todos []model.Todo
	for rows.Next() {
		var t model.Todo
		if err := rows.Scan(&t.ID, &t.EventID, &t.Content, &t.Done); err != nil {
			return nil, fmt.Errorf("scan todo: %w", err)
		}
		todos = append(todos, t)
	}
	cal.Todos = todos

	return &cal, nil
}

func (r *CalendarRepository) DeleteCalendar(ctx context.Context, eventID uint64) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// 1️⃣ Todos 먼저 삭제
	if _, err := tx.ExecContext(ctx, `DELETE FROM todos WHERE event_id = ?`, eventID); err != nil {
		tx.Rollback()
		return fmt.Errorf("delete todos: %w", err)
	}

	// 2️⃣ Calendar 삭제
	if _, err := tx.ExecContext(ctx, `DELETE FROM calendars WHERE event_id = ?`, eventID); err != nil {
		tx.Rollback()
		return fmt.Errorf("delete calendar: %w", err)
	}

	// 3️⃣ Commit
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("transaction commit failed: %w", err)
	}

	return nil
}
