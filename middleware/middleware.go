package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rainbow96bear/planet_user_server/internal/service"
	"github.com/rainbow96bear/planet_user_server/utils"
	"github.com/rainbow96bear/planet_utils/pkg/logger"
)

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		method := c.Request.Method

		logger.Infof("start request: %s %s", method, path)
		start := time.Now()

		// 다음 핸들러 실행
		c.Next()

		duration := time.Since(start)
		logger.Infof("end request: %s %s, status=%d, duration=%s",
			method, path, c.Writer.Status(), duration)
	}
}

func AuthMiddleware(authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 쿠키에서 access_token 가져오기
		tokenStr, err := c.Cookie("access_token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing access token"})
			c.Abort()
			return
		}

		userUuid, err := utils.GetUuidByAccessToken(tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		nickname := c.Param("nickname")
		if nickname == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing nickname param"})
			c.Abort()
			return
		}

		ctx := c.Request.Context()
		// 두 값을 비교
		ok, err := authService.VerifyUser(ctx, nickname, userUuid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized user"})
			c.Abort()
			return
		}

		// 검증 성공 시 userUuid를 Context에 저장
		c.Set("userUuid", userUuid)
		c.Next()
	}
}
