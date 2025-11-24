package service

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/rainbow96bear/planet_utils/models"
)

type CalendarCacheItem struct {
	Data      []*models.CalendarEvents // DB 모델 그대로 캐시
	ExpiresAt time.Time
}

// key: "UserID:year:month:visibility"
var calendarCache sync.Map
var cacheTTL = 1 * time.Minute

// BuildCacheKey generates key for caching
func buildCacheKey(UserID uuid.UUID, year int, month int, visibility string) string {
	UserIDStr := UserID.String()
	return UserIDStr + ":" +
		time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC).Format("2006-01") +
		":" + visibility
}

// GetCalendarCache retrieves cached calendars if exists and not expired
func GetCalendarCache(UserID uuid.UUID, year int, month int, visibility string) ([]*models.CalendarEvents, bool) {
	key := buildCacheKey(UserID, year, month, visibility)
	if item, ok := calendarCache.Load(key); ok {
		cacheItem, valid := item.(CalendarCacheItem)
		if valid && time.Now().Before(cacheItem.ExpiresAt) {
			return cacheItem.Data, true
		}
		calendarCache.Delete(key)
	}
	return nil, false
}

// SetCalendarCache stores calendar data in cache
func SetCalendarCache(UserID uuid.UUID, year int, month int, visibility string, data []*models.CalendarEvents) {
	key := buildCacheKey(UserID, year, month, visibility)
	calendarCache.Store(key, CalendarCacheItem{
		Data:      data,
		ExpiresAt: time.Now().Add(cacheTTL),
	})
}

func DeleteCalendarCache(UserID uuid.UUID, year int, month int, visibility string) {
	key := buildCacheKey(UserID, year, month, visibility)
	calendarCache.Delete(key)
}
