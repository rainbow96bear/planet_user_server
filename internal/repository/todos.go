package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/rainbow96bear/planet_user_server/internal/models"
	"github.com/rainbow96bear/planet_user_server/internal/tx"
	"github.com/rainbow96bear/planet_utils/pkg/logger"
	"gorm.io/gorm"
)

type TodosRepository struct {
	db *gorm.DB
}

func NewTodosRepository(db *gorm.DB) *TodosRepository {
	if db == nil {
		panic("database connection is required")
	}
	return &TodosRepository{
		db: db,
	}
}

func (r *TodosRepository) getDB(ctx context.Context) *gorm.DB {
	// tx 패키지를 사용하여 Context에서 트랜잭션을 추출합니다.
	if tx := tx.GetTx(ctx); tx != nil {
		return tx.WithContext(ctx)
	}
	return r.db.WithContext(ctx) // 기본 DB 연결 반환
}

// func (r *TodosRepository) BeginTx(ctx context.Context) (*gorm.DB, error) {
// 	logger.Infof("starting transaction for TodosRepository")
// 	tx := r.DB.WithContext(ctx).Begin()
// 	if tx.Error != nil {
// 		logger.Errorf("failed to start transaction: %v", tx.Error)
// 		return nil, tx.Error
// 	}
// 	logger.Infof("transaction started successfully")
// 	return tx, nil
// }

// -------------------------
// 단일 Todo 생성
// -------------------------
func (r *TodosRepository) CreateTodos(
	ctx context.Context,
	todos []models.Todo,
) error {
	if len(todos) == 0 {
		logger.Debugf("[TodosRepo] no todos to create")
		return nil
	}

	db := r.getDB(ctx)

	logger.Infof(
		"[TodosRepo] creating %d todos (calendar_event_id=%s)",
		len(todos),
		todos[0].CalendarEventID,
	)

	if err := db.WithContext(ctx).Create(&todos).Error; err != nil {
		logger.Errorf(
			"[TodosRepo] failed to create todos (calendar_event_id=%s): %v",
			todos[0].CalendarEventID,
			err,
		)
		return fmt.Errorf("failed to create todos: %w", err)
	}

	logger.Infof(
		"[TodosRepo] successfully created %d todos (calendar_event_id=%s)",
		len(todos),
		todos[0].CalendarEventID,
	)

	return nil
}

// -------------------------
// Todo 상태 업데이트
// -------------------------
func (r *TodosRepository) UpdateTodoStatus(
	ctx context.Context,
	userID uuid.UUID,
	todoID uuid.UUID,
	isDone bool,
) (*models.Todo, error) {

	db := r.getDB(ctx)

	var todo models.Todo

	// 1️⃣ Todo + Event 조인해서 소유권 확인
	if err := db.
		Joins("JOIN calendar_events ce ON ce.id = todos.calendar_event_id").
		Where("todos.id = ? AND ce.user_id = ?", todoID, userID).
		First(&todo).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("unauthorized or todo not found")
		}
		return nil, err
	}

	// 2️⃣ 상태 변경
	todo.IsDone = isDone
	if err := db.Save(&todo).Error; err != nil {
		return nil, err
	}

	return &todo, nil
}

func (r *TodosRepository) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (*models.Todo, error) {
	db := r.getDB(ctx)
	var todo models.Todo
	err := db.WithContext(ctx).
		Preload("CalendarEvent").
		First(&todo, "id = ?", id).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &todo, nil
}

// // -------------------------
// // EventID 기반 Todo 조회
// // -------------------------
// func (r *TodosRepository) FindTodosByEventID(ctx context.Context, eventID uuid.UUID) ([]*models.Todos, error) {
// 	var todos []*models.Todos
// 	if err := r.DB.WithContext(ctx).Where("event_id = ?", eventID).Find(&todos).Error; err != nil {
// 		return nil, fmt.Errorf("failed to fetch todos: %w", err)
// 	}
// 	logger.Infof("Found %d todos for event %s", len(todos), eventID)
// 	return todos, nil
// }

// // -------------------------
// // 단일 Todo 삭제
// // -------------------------
// func (r *TodosRepository) DeleteTodo(ctx context.Context, todoID uuid.UUID) error {
// 	if err := r.DB.WithContext(ctx).Delete(&models.Todos{}, "id = ?", todoID).Error; err != nil {
// 		return fmt.Errorf("failed to delete todo: %w", err)
// 	}
// 	logger.Infof("Todo %s deleted", todoID)
// 	return nil
// }
