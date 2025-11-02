package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rainbow96bear/planet_user_server/config"
	"github.com/rainbow96bear/planet_utils/pkg/jwt"
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

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing Authorization header"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid Authorization header format"})
			c.Abort()
			return
		}

		tokenStr := parts[1]

		// JWT 검증
		claims, err := jwt.ParseAndVerifyJWT(tokenStr, config.JWT_SECRET_KEY)
		logger.Debugf("tokenStr : %s", tokenStr)
		logger.Debugf("claims : %v", claims)
		if err != nil {
			logger.Errorf("invalid token: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			c.Abort()
			return
		}

		// 필수 claim 확인
		userUuid, ok1 := claims["user_uuid"].(string)
		nickname, ok2 := claims["nickname"].(string)
		if !ok1 || !ok2 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			c.Abort()
			return
		}

		// Context에 저장 (핸들러에서 사용 가능)
		c.Set("user_uuid", userUuid)
		c.Set("nickname", nickname)

		c.Next()
	}
}
