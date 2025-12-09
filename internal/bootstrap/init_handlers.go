package bootstrap

import (

	// 프로젝트 내부 패키지

	"github.com/rainbow96bear/planet_user_server/internal/handler"
)

// HandlerMap: 모든 초기화된 핸들러를 저장하는 맵
type HandlerMap map[string]handler.Handler

// InitDependencies: 모든 하위 의존성(Repo, Client, Service)을 초기화하고 HTTP 핸들러를 반환합니다.
// 이 함수는 'main'에서 호출됩니다.
func InitHandlers(dep *Dependencies) HandlerMap {
	graphqlHandler := handler.NewGraphqlHandler(dep.Resolver)

	return HandlerMap{
		"graphql": graphqlHandler,
	}
}
