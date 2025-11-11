package service

import (
	"sync"
	"time"

	"github.com/rainbow96bear/planet_utils/model"
)

type CalendarCacheItem struct {
	Data      []*model.Calendar // DB 모델 그대로 캐시
	ExpiresAt time.Time
}

// key: "userUUID:year:month:visibility"
var calendarCache sync.Map
var cacheTTL = 1 * time.Minute

// BuildCacheKey generates key for caching
func buildCacheKey(userUUID string, year int, month int, visibility string) string {
	return userUUID + ":" +
		time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC).Format("2006-01") +
		":" + visibility
}

// GetCalendarCache retrieves cached calendars if exists and not expired
func GetCalendarCache(userUUID string, year int, month int, visibility string) ([]*model.Calendar, bool) {
	key := buildCacheKey(userUUID, year, month, visibility)
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
func SetCalendarCache(userUUID string, year int, month int, visibility string, data []*model.Calendar) {
	key := buildCacheKey(userUUID, year, month, visibility)
	calendarCache.Store(key, CalendarCacheItem{
		Data:      data,
		ExpiresAt: time.Now().Add(cacheTTL),
	})
}

func DeleteCalendarCache(userUUID string, year int, month int, visibility string) {
	key := buildCacheKey(userUUID, year, month, visibility)
	calendarCache.Delete(key)
}
