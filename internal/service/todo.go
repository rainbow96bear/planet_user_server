package service

// TodoService: Todo 항목 관리를 전담합니다.
// type TodoService struct {
// 	// CalendarEventsRepo를 통해 Todo 테이블에 접근합니다.
// 	TodosRepo *repository.TodosRepository
// }

// // NewTodoService: TodoService를 생성합니다.
// func NewTodoService(todosRepo *repository.TodosRepository) *TodoService {
// 	return &TodoService{
// 		TodosRepo: todosRepo,
// 	}
// }

// // ----------------------------
// // Todo 상태 업데이트
// // ----------------------------

// // UpdateTodoStatus: 특정 Todo 항목의 isDone 상태를 업데이트하고, 관련된 Event 캐시를 무효화합니다.
// // 💡 이 함수는 Handler에서 직접 호출됩니다.
// func (s *TodoService) UpdateTodoStatus(ctx context.Context, userID uuid.UUID, todoID uuid.UUID, isDone bool) error {
// 	logger.Infof("[TodoService.UpdateTodoStatus] UserID=%s, TodoID=%s, IsDone=%t", userID, todoID, isDone)

// 	// 1. Repository를 통해 Todo 상태 업데이트 및 소유권 확인
// 	// Repository는 업데이트 성공 시 해당 Todo가 속한 Event 정보를 반환해야 합니다.
// 	err := s.TodosRepo.UpdateTodoStatus(ctx, todoID, isDone)

// 	if err != nil {
// 		// 예: unauthorized, not found 등의 에러를 Repository에서 반환한다고 가정합니다.
// 		logger.Errorf("[TodoService.UpdateTodoStatus] failed to update todo: %v", err)
// 		return err
// 	}

// 	logger.Infof("[TodoService.UpdateTodoStatus] Todo status updated successfully: %s", todoID)
// 	return nil
// }

// // ----------------------------
// // (추가 예정) 기타 Todo 관련 CRUD (예: Todo 개별 생성/수정/삭제)
// // ----------------------------
// // func (s *TodoService) DeleteTodo(ctx context.Context, userID uuid.UUID, todoID uuid.UUID) error { ... }
// // func (s *TodoService) FindTodoByID(ctx context.Context, todoID uuid.UUID) (*dto.TodoItem, error) { ... }
