// middleware/auth_middleware.go (Gin 전용으로 수정)

package middleware

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const ContextKeyAccessToken contextKey = "access_token"

var (
	ErrAuthorizationMissing = errors.New("authorization header missing")
	ErrMalformedToken       = errors.New("malformed authorization token")
)

// 🚨 Gin 전용 미들웨어 함수: gin.HandlerFunc를 반환합니다.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		token := ""

		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			token = strings.TrimPrefix(authHeader, "Bearer ")
		}

		// Context를 복사하고 토큰 값을 주입
		ctx := context.WithValue(c.Request.Context(), ContextKeyAccessToken, token)

		// 업데이트된 Context로 요청 객체 대체
		c.Request = c.Request.WithContext(ctx)

		// 다음 핸들러(GraphQL 서버)로 체인 전달
		c.Next()
	}
}

func ExtractAccessToken(ctx context.Context) (*jwt.Token, error) {
	accessToken, ok := ctx.Value(ContextKeyAccessToken).(string)

	if !ok || accessToken == "" {
		return nil, ErrAuthorizationMissing
	}

	tokenString := strings.TrimPrefix(accessToken, "Bearer ")

	if tokenString == "" {
		tokenString = accessToken
	}

	if tokenString == "" {
		return nil, ErrMalformedToken
	}
	token, _, err := jwt.NewParser().ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}
	return token, nil
}
