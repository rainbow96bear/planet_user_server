package handler

// import (
// 	"fmt"
// 	"net/http"

// 	// time 패키지는 현재 사용하지 않지만, CalendarService에 필요할 수 있으므로 유지
// 	"github.com/gin-gonic/gin"
// 	"github.com/google/uuid"
// 	"github.com/rainbow96bear/planet_user_server/dto"
// 	"github.com/rainbow96bear/planet_user_server/internal/service"
// 	"github.com/rainbow96bear/planet_user_server/middleware"
// 	"github.com/rainbow96bear/planet_user_server/utils"
// 	"github.com/rainbow96bear/planet_utils/pkg/logger"
// )

// type TodoHandler struct {
// 	TodoService *service.TodoService
// }

// func NewTodoHandler(todoService *service.TodoService) *TodoHandler {
// 	return &TodoHandler{TodoService: todoService}
// }

// func (h *TodoHandler) RegisterRoutes(r *gin.Engine) {
// 	// PlanHandler의 /me 그룹을 사용하여 인증 미들웨어 적용
// 	me := r.Group("/me")
// 	me.Use(middleware.AccessTokenAuthMiddleware())
// 	{
// 		// 💡 To-do 리소스 그룹 정의
// 		todos := me.Group("/todos")
// 		{
// 			// PATCH /me/todos/:todoId
// 			// todoId에 해당하는 항목의 is_done 상태를 업데이트합니다.
// 			todos.PATCH("/:todoId", h.UpdateTodoStatus)
// 		}
// 	}
// }

// // ---------------------- Handler Implementations ----------------------

// // UpdateTodoStatus: 특정 Todo 항목의 is_done 상태를 업데이트합니다.
// func (h *TodoHandler) UpdateTodoStatus(c *gin.Context) {
// 	ctx := c.Request.Context()

// 	// 1. User ID 추출 (Auth Middleware에서 설정)
// 	userID, err := utils.GetUserID(c)
// 	if err != nil {
// 		logger.Errorf("UpdateTodoStatus failed: UserID not found in context")
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "user ID not found"})
// 		return
// 	}

// 	// 2. todoId 파라미터 추출 및 유효성 검사
// 	todoIDStr := c.Param("todoId")
// 	todoID, err := uuid.Parse(todoIDStr)
// 	if err != nil {
// 		logger.Warnf("UpdateTodoStatus received invalid todoId: %s", todoIDStr)
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid todo ID format"})
// 		return
// 	}

// 	// 3. 요청 본문 파싱
// 	var req dto.TodoUpdateStatusRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		logger.Warnf("UpdateTodoStatus failed binding JSON for todoID=%s: %v", todoID, err)
// 		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid request body: %v", err)})
// 		return
// 	}

// 	logger.Infof("UpdateTodoStatus start: UserID=%s, TodoID=%s, IsDone=%t", userID, todoID, req.IsDone)

// 	// 4. Service 로직 호출
// 	err = h.TodoService.UpdateTodoStatus(ctx, userID, todoID, req.IsDone)

// 	if err != nil {
// 		logger.Errorf("UpdateTodoStatus failed for todoID=%s: %v", todoID, err)

// 		// 에러 타입에 따라 세분화된 응답 제공 (예: 권한 없음, 찾을 수 없음)
// 		if err.Error() == "unauthorized" {
// 			c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized access to todo"})
// 			return
// 		}
// 		if err.Error() == "not found" {
// 			c.JSON(http.StatusNotFound, gin.H{"error": "todo item not found"})
// 			return
// 		}

// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update todo status"})
// 		return
// 	}

// 	logger.Infof("UpdateTodoStatus successful: TodoID=%s", todoID)
// 	c.Status(http.StatusNoContent) // 성공적인 업데이트 후 본문 없음 응답 (204 No Content)
// }
