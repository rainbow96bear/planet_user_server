package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rainbow96bear/planet_user_server/dto"
	"github.com/rainbow96bear/planet_user_server/internal/service"
	"github.com/rainbow96bear/planet_user_server/middleware"
	"github.com/rainbow96bear/planet_user_server/utils"
	"github.com/rainbow96bear/planet_utils/pkg/logger"
)

type CalendarHandler struct {
	CalendarService *service.CalendarService
	FollowService   *service.FollowService
	ProfileService  *service.ProfileService
	// 추후 이미지 업로드 서비스 추가
	// ImageUploader dto.ImageUploader
}

func NewCalendarHandler(calendarService *service.CalendarService) *CalendarHandler {
	return &CalendarHandler{
		CalendarService: calendarService,
	}
}

func (h *CalendarHandler) RegisterRoutes(r *gin.Engine) {
	calendarGroup := r.Group("/calendar")

	// 공개용 (다른 사람 일정 조회)
	calendarGroup.GET("/user/:nickname", h.GetUserCalendar)

	// 인증 필요
	calendarGroup.Use(middleware.AuthMiddleware())
	{
		calendarGroup.GET("", h.GetMyCalendar)
		calendarGroup.POST("", h.CreateCalendar)
		calendarGroup.PUT("/:eventId", h.UpdateCalendar)
		calendarGroup.DELETE("/:eventId", h.DeleteCalendar)
	}
}

func (h *CalendarHandler) GetUserCalendar(c *gin.Context) {
	logger.Infof("start to get user calendar")
	defer logger.Infof("end to get user calendar")

	ctx := c.Request.Context()
	nickname := c.Param("nickname")
	if nickname == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "nickname is required"})
		return
	}

	// year/month 쿼리
	yearStr := c.Query("year")
	monthStr := c.Query("month")
	if yearStr == "" || monthStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "year and month are required"})
		return
	}
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid year"})
		return
	}
	month, err := strconv.Atoi(monthStr)
	if err != nil || month < 1 || month > 12 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid month"})
		return
	}

	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0) // 다음 달 1일

	// visibility
	authUuid, _ := utils.GetUserUuid(c)

	visibilityLevel := []string{"public"}
	if authUuid != "" {
		followeeUuid, err := h.ProfileService.GetUserUuidByNickname(ctx, nickname)
		if err == nil {
			isFollow, err := h.FollowService.IsFollow(ctx, authUuid, followeeUuid)
			if err == nil && isFollow {
				visibilityLevel = append(visibilityLevel, "friends")
			}
		}
	}

	// DB 조회 + 캐시
	calendars, err := h.CalendarService.GetUserCalendars(ctx, nickname, visibilityLevel, startDate, endDate)
	if err != nil {
		logger.Errorf("failed to get calendars: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get calendars"})
		return
	}

	eventsResp := dto.ToCalendarInfoList(calendars)
	monthData := h.CalendarService.GenerateMonthData(startDate)
	completionData := h.CalendarService.CalculateCompletionData(calendars)

	c.JSON(http.StatusOK, gin.H{
		"events":         eventsResp,
		"monthData":      monthData,
		"completionData": completionData,
	})
}

func (h *CalendarHandler) GetMyCalendar(c *gin.Context) {
	ctx := c.Request.Context()
	userUUID, _ := utils.GetUserUuid(c) // AuthMiddleware에서 세팅

	// year/month 쿼리
	yearStr := c.Query("year")
	monthStr := c.Query("month")
	if yearStr == "" || monthStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "year and month are required"})
		return
	}
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid year"})
		return
	}
	month, err := strconv.Atoi(monthStr)
	if err != nil || month < 1 || month > 12 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid month"})
		return
	}

	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0) // 다음 달 1일

	// 내 캘린더는 private 포함
	visibility := []string{"public", "friends", "private"}

	calendars, err := h.CalendarService.GetUserCalendars(ctx, userUUID, visibility, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get calendars"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"events":         dto.ToCalendarInfoList(calendars),
		"monthData":      h.CalendarService.GenerateMonthData(startDate),
		"completionData": h.CalendarService.CalculateCompletionData(calendars),
	})
}

func (h *CalendarHandler) CreateCalendar(c *gin.Context) {
	ctx := c.Request.Context()
	userUUID, err := utils.GetUserUuid(c)
	if err != nil {
		logger.Errorf(err.Error())
	}

	var req dto.CalendarCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "detail": err.Error()})
		return
	}

	// DTO -> Model 변환
	calendar := dto.ToCalendarModelFromCreate(&req, userUUID)

	// DB 저장
	if err := h.CalendarService.CreateCalendar(ctx, calendar); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create calendar"})
		return
	}

	// 캐시 갱신
	h.CalendarService.ClearCache(userUUID, calendar.StartAt.Year(), int(calendar.StartAt.Month()))

	c.JSON(http.StatusCreated, gin.H{
		"message": "calendar created",
		"event":   dto.ToCalendarInfo(calendar),
	})
}

func (h *CalendarHandler) UpdateCalendar(c *gin.Context) {}

func (h *CalendarHandler) DeleteCalendar(c *gin.Context) {
	ctx := c.Request.Context()

	// eventId 파라미터 확인
	eventIDStr := c.Param("eventId")
	if eventIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "eventId is required"})
		return
	}
	eventID, err := strconv.ParseUint(eventIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid eventId"})
		return
	}

	// 인증된 유저 UUID 가져오기
	userUUID, err := utils.GetUserUuid(c)
	if err != nil || userUUID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// DB에서 삭제 (Repository/Service 호출)
	if err := h.CalendarService.DeleteCalendar(ctx, userUUID, eventID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete calendar"})
		return
	}

	// 캐시 삭제 (삭제된 이벤트의 연도/월/visibility 기준)
	// Service 내부에서 DeleteCalendarCache를 호출하도록 처리하면 좋음
	// 예: h.CalendarService.DeleteCalendarCache(userUUID, year, month, visibility)

	c.JSON(http.StatusOK, gin.H{"message": "calendar deleted successfully"})
}
