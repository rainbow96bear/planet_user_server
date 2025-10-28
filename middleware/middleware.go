package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rainbow96bear/planet_auth_server/config"
	"github.com/rainbow96bear/planet_user_server/internal/service"
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

func AuthMiddleware(authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1️⃣ Authorization 헤더에서 access token 가져오기
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

		// 2️⃣ JWT 검증 및 payload 추출
		claims, err := jwt.ParseAndVerifyJWT(tokenStr, config.JWT_SECRET_KEY)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			c.Abort()
			return
		}

		// claims에서 userUuid와 nickname 가져오기
		userUuid, ok1 := claims["userUuid"].(string)
		nicknameFromToken, ok2 := claims["nickname"].(string)
		if !ok1 || !ok2 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			c.Abort()
			return
		}

		// 3️⃣ 요청 URL의 nickname과 JWT nickname 비교
		nickname := c.Param("nickname")
		if nickname == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing nickname param"})
			c.Abort()
			return
		}

		if nickname != nicknameFromToken {
			c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized user"})
			c.Abort()
			return
		}

		// 검증 성공 → Context에 저장
		c.Set("userUuid", userUuid)
		c.Set("nickname", nicknameFromToken)
		c.Next()
	}
}
