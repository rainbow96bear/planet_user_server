package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rainbow96bear/planet_user_server/internal/service"
	"github.com/rainbow96bear/planet_user_server/middleware"
	"github.com/rainbow96bear/planet_user_server/utils"
	"github.com/rainbow96bear/planet_utils/pkg/logger"
)

type PlanHandler struct {
	CalendarService *service.CalendarService
	// TodoService *service.TodoService // (만약 TodoService가 분리된다면 추가)
}

func NewPlanHandler(calendarService *service.CalendarService) *PlanHandler {
	return &PlanHandler{
		CalendarService: calendarService,
	}
}

func (h *PlanHandler) RegisterRoutes(r *gin.Engine) {
	// PlanHandler는 /plans API 그룹을 담당합니다.
	me := r.Group("/me")
	me.Use(middleware.AccessTokenAuthMiddleware())
	{
		// 1. 내 일일 계획 조회 (Event + Nested Todos 포함)
		// GET /me/plans/daily?date=YYYY-MM-DD
		me.GET("/plans/daily", h.GetMyDailyPlan)
	}

	users := r.Group("/users/:nickname")
	{
		// 2. 다른 사용자 일일 계획 조회 (Event + Nested Todos 포함)
		// GET /users/:nickname/plans/daily?date=YYYY-MM-DD
		users.GET("/plans/daily", h.GetUserDailyPlan)
	}
}

// ---------------------- Handler Implementations ----------------------

// GetMyDailyPlan: 내 일일 계획 (Event + Todo) 조회
func (h *PlanHandler) GetMyDailyPlan(c *gin.Context) {
	ctx := c.Request.Context()
	// utils.GetUserID는 UserID와 에러를 반환해야 하지만,
	// AccessTokenAuthMiddleware 덕분에 UserID가 존재한다고 가정하고 에러는 무시합니다.
	UserID, _ := utils.GetUserID(c)
	logger.Infof("GetMyDailyPlan start, UserID=%s", UserID)
	defer logger.Infof("GetMyDailyPlan end, UserID=%s", UserID)

	// 'date' 쿼리 파라미터 파싱
	date, err := parseDate(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// CalendarService를 통해 Event와 Todo 데이터를 통합하여 조회
	// 이 메서드(GetMyCalendarDailyData)는 CalendarService 내부에 구현되어야 합니다.
	data, err := h.CalendarService.GetMyCalendarDailyData(ctx, UserID, date)
	if err != nil {
		logger.Errorf("GetMyDailyPlan failed, UserID=%s, date=%s, err=%v", UserID, date.Format("2006-01-02"), err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get daily plans"})
		return
	}

	c.JSON(http.StatusOK, data)
}

// GetUserDailyPlan: 다른 사용자 일일 계획 (Event + Todo) 조회
func (h *PlanHandler) GetUserDailyPlan(c *gin.Context) {
	ctx := c.Request.Context()
	nickname := c.Param("nickname")
	logger.Infof("GetUserDailyPlan start, nickname=%s", nickname)
	defer logger.Infof("GetUserDailyPlan end, nickname=%s", nickname)

	if nickname == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "nickname is required"})
		return
	}

	// 'date' 쿼리 파라미터 파싱
	date, err := parseDate(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 인증된 사용자 ID (공개 범위 확인용)
	authID, _ := utils.GetUserID(c)

	// CalendarService를 통해 조회 (권한 확인 로직은 서비스 레이어에 구현되어야 함)
	data, err := h.CalendarService.GetUserCalendarDailyData(ctx, nickname, authID, date)
	if err != nil {
		logger.Errorf("GetUserDailyPlan failed, nickname=%s, date=%s, err=%v", nickname, date.Format("2006-01-02"), err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get daily plans"})
		return
	}

	c.JSON(http.StatusOK, data)
}

// ---------------------- Helper ----------------------
// parseDate: 'date' 쿼리 파라미터를 YYYY-MM-DD 형식으로 파싱합니다.
// 이 함수는 utils나 handler 공통 파일에 정의되어 있어야 합니다.
func parseDate(c *gin.Context) (time.Time, error) {
	dateStr := c.Query("date")
	if dateStr == "" {
		return time.Time{}, fmt.Errorf("date query parameter is required in YYYY-MM-DD format")
	}

	// YYYY-MM-DD 형식 파싱
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date format, must be YYYY-MM-DD")
	}

	return date, nil
}
