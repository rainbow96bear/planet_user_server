// 파일: main.go (개선 및 통합)

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"

	// 🌟 프로젝트 내부 패키지
	"github.com/rainbow96bear/planet_user_server/config"
	"github.com/rainbow96bear/planet_user_server/internal/bootstrap"
	grpcserver "github.com/rainbow96bear/planet_user_server/internal/grpc/server"
	"github.com/rainbow96bear/planet_user_server/internal/router"

	// 🌟 공통 유틸리티/프로토 버퍼

	"github.com/rainbow96bear/planet_utils/pkg/logger"
)

// 빌드 플래그 (생략)
var (
	Mode      string
	Version   string
	GitCommit string
)

func init() {
	versionFlag := flag.Bool("version", false, "Print version and exit")
	flag.Parse()

	if *versionFlag {
		fmt.Printf("Version: %s\nCommit: %s\n", Version, GitCommit)
		os.Exit(0)
	}
	Mode = "dev"
	fmt.Printf("user_server Start \nVersion : %s \nGit Commit : %s\n", Version, GitCommit)
	fmt.Printf("Build Mode : %s\n", Mode)
	config.InitConfig(Mode)
	logger.SetLevel(config.LOG_LEVEL)
}

func main() {
	// ----------------------------------------------------------------------
	// 0. 인프라 초기화 (DB 연결)
	// ----------------------------------------------------------------------
	db, err := bootstrap.InitDatabase()
	if err != nil {
		logger.Errorf("failed to initialize database: %v", err)
		os.Exit(1)
	}
	sqlDB, err := db.DB()

	// 2. 에러 핸들링 및 defer를 사용하여 프로그램 종료 시 연결을 닫습니다.
	if err != nil {
		// 내부 sql.DB를 얻는 데 실패하면 경고만 로깅합니다.
		logger.Warnf("failed to get underlying sql.DB for closing: %v", err)
	} else {
		// 💡 defer를 사용하여 main 함수 종료 시 연결 풀을 안전하게 닫습니다.
		defer func() {
			if closeErr := sqlDB.Close(); closeErr != nil {
				logger.Errorf("failed to close database connection: %v", closeErr)
			}
		}()
	}

	dependencies, err := bootstrap.InitDependencies(db)
	if err != nil {
		logger.Errorf("fail to init Dependencies %s", err.Error())
		os.Exit(1)
	}

	go grpcserver.RunGrpcServer(db, dependencies)

	// ----------------------------------------------------------------------
	// HTTP/GraphQL 서버 실행 (Gin)
	// ----------------------------------------------------------------------

	// 💡 컨테이너에서 UserService를 꺼내 GraphQL Resolver에 주입합니다.
	// GraphQL Resolver는 DB 대신 Service 계층에 의존해야 합니다.

	handlers := bootstrap.InitHandlers(dependencies)

	r := router.SetupRouter(func(r *gin.Engine) {
		for _, h := range handlers {
			h.RegisterRoutes(r)
		}
	})

	userServerPort := fmt.Sprintf(":%s", config.PORT)
	logger.Infof("GraphQL/HTTP Server started on port %s", config.PORT)

	if err := r.Run(userServerPort); err != nil {
		logger.Errorf("failed to start http server: %v", err)
		os.Exit(1)
	}
}
