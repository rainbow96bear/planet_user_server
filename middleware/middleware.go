package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	pb "github.com/rainbow96bear/planet_proto"
	"github.com/rainbow96bear/planet_user_server/config"
	"github.com/rainbow96bear/planet_user_server/grpc_client"
	"github.com/rainbow96bear/planet_user_server/logger"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 🍪 1. 쿠키에서 access_token 가져오기
		tokenStr, err := c.Cookie("access_token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing access token"})
			c.Abort()
			return
		}

		// 🔐 2. 토큰 검증
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.JWT_SECRET_KEY), nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid access token"})
			c.Abort()
			return
		}

		// 🧩 3. 클레임에서 userUuid 추출
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			c.Abort()
			return
		}

		userUuid, ok := claims["user_uuid"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user uuid in token"})
			c.Abort()
			return
		}

		// 🧭 4. params에서 nickname 추출
		nickname := c.Param("nickname")
		if nickname == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing nickname param"})
			c.Abort()
			return
		}

		// 🗄️ 5. nickname으로 DB 조회해서 uuid 확인
		dbClient, err := grpc_client.NewDBClient(config.DB_GRPC_SERVER_ADDR)
		if err != nil {
			logger.Errorf("fail to create dbClient ERR[%s]", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify user"})
			c.Abort()
			return
		}

		reqUserInfo := &pb.UserInfo{
			Nickname: nickname,
		}

		userInfo, err := dbClient.ReqGetUserInfoByNickname(reqUserInfo)
		if err != nil {
			logger.Errorf("failed to query user by nickname ERR[%s]", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify user"})
			c.Abort()
			return
		}

		// ⚖️ 6. JWT의 userUuid와 DB의 userUuid 비교
		if userInfo.UserUuid != userUuid {
			c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized user"})
			c.Abort()
			return
		}

		// ✅ 검증 성공 시 userUuid를 Context에 저장
		c.Set("userUuid", userUuid)
		c.Next()
	}
}
