package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/rainbow96bear/planet_user_server/internal/models"
	"github.com/rainbow96bear/planet_user_server/internal/repository"
	"github.com/rainbow96bear/planet_utils/pkg/logger"
	"gorm.io/gorm"
)

// TodoService: Todo 항목 관리를 전담합니다.
type TodoServiceInterface interface {
	UpdateTodoStatus(
		ctx context.Context,
		userID uuid.UUID,
		todoID uuid.UUID,
		isDone bool,
	) (*models.Todo, error)
	FindByID(
		ctx context.Context,
		userID uuid.UUID,
		todoID uuid.UUID,
	) (*models.Todo, error)
}

type TodoService struct {
	db *gorm.DB
	// CalendarEventsRepo를 통해 Todo 테이블에 접근합니다.
	TodosRepo *repository.TodosRepository
}

// NewTodoService: TodoService를 생성합니다.
func NewTodoService(db *gorm.DB, todosRepo *repository.TodosRepository) *TodoService {
	return &TodoService{
		db:        db,
		TodosRepo: todosRepo,
	}
}

// // ----------------------------
// // Todo 상태 업데이트
// // ----------------------------

// UpdateTodoStatus: 특정 Todo 항목의 isDone 상태를 업데이트하고, 관련된 Event 캐시를 무효화합니다.
// 💡 이 함수는 Handler에서 직접 호출됩니다.
func (s *TodoService) UpdateTodoStatus(
	ctx context.Context,
	userID uuid.UUID,
	todoID uuid.UUID,
	isDone bool,
) (*models.Todo, error) {

	logger.Infof(
		"[TodoService.UpdateTodoStatus] user=%s todo=%s done=%t",
		userID, todoID, isDone,
	)

	// Repository에서 권한 검증 + 업데이트 + 반환까지
	todo, err := s.TodosRepo.UpdateTodoStatus(
		ctx,
		userID,
		todoID,
		isDone,
	)
	if err != nil {
		return nil, err
	}

	return todo, nil
}

// // ----------------------------
// // (추가 예정) 기타 Todo 관련 CRUD (예: Todo 개별 생성/수정/삭제)
// // ----------------------------
// // func (s *TodoService) DeleteTodo(ctx context.Context, userID uuid.UUID, todoID uuid.UUID) error { ... }
func (s *TodoService) FindByID(
	ctx context.Context,
	userID uuid.UUID,
	todoID uuid.UUID,
) (*models.Todo, error) {

	logger.Infof(
		"[TodoService.FindByID] UserID=%s TodoID=%s",
		userID, todoID,
	)

	todo, err := s.TodosRepo.FindByID(ctx, todoID)
	if err != nil {
		return nil, err
	}
	if todo == nil {
		return nil, fmt.Errorf("todo not found")
	}

	return todo, nil
}
