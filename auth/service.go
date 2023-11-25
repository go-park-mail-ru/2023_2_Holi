package auth

import (
	csrf_http "2023_2_Holi/csrf/delivery/http"
	"2023_2_Holi/domain"
	"context"
	"net"
	"net/http"
	"os"

	"google.golang.org/grpc"

	auth_grpc "2023_2_Holi/auth/delivery/grpc"
	auth_http "2023_2_Holi/auth/delivery/http"
	auth_postgres "2023_2_Holi/auth/repository/postgresql"
	auth_redis "2023_2_Holi/auth/repository/redis"
	auth_usecase "2023_2_Holi/auth/usecase"
	"2023_2_Holi/connectors/postgres"
	"2023_2_Holi/connectors/redis"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"2023_2_Holi/domain/grpc/session"
	logs "2023_2_Holi/logger"
	"2023_2_Holi/middleware"
)

func startRpcServer(au domain.AuthUsecase) {
	port := ":" + os.Getenv("AUTHMS_GRPC_SERVER_PORT")
	lis, err := net.Listen("tcp", port)
	if err != nil {
		logs.LogFatal(logs.Logger, "auth", "startRpcServer", err, err.Error())
	}

	server := grpc.NewServer()

	session.RegisterAuthCheckerServer(server, auth_grpc.NewAuthAuthHandler(au))

	logs.Logger.Info("starting auth grpc server at ", port)
	server.Serve(lis)
	logs.Logger.Info("auth grpc server stopped")
}

func StartService() {
	err := godotenv.Load()
	ctx := context.Background()
	accessLogger := middleware.AccessLogger{
		LogrusLogger: logs.Logger,
	}

	dbParams := postgres.GetParamsForNetflixDB()
	pc := postgres.Connect(ctx, dbParams)
	defer pc.Close()
	rc := redis.Connect()
	defer rc.Close()

	//authMiddlewareRouter := mainRouter.PathPrefix("/api").Subrouter()

	sr := auth_redis.NewSessionRedisRepository(rc)
	ar := auth_postgres.NewAuthPostgresqlRepository(pc, ctx)

	au := auth_usecase.NewAuthUsecase(ar, sr)

	mainRouter := mux.NewRouter()
	tokens, _ := domain.NewHMACHashToken("Gvjhlk123bl1lma0")

	auth_http.NewAuthHandler(mainRouter, au)
	csrf_http.NewCsrfHandler(mainRouter, tokens)

	go startRpcServer(au)

	mw := middleware.InitMiddleware(nil, tokens)

	//authMiddlewareRouter.Use(mw.IsAuth)
	mainRouter.Use(accessLogger.AccessLogMiddleware)
	mainRouter.Use(mux.CORSMethodMiddleware(mainRouter))
	mainRouter.Use(mw.CORS)
	mainRouter.Use(mw.CSRFProtection)
	mainRouter.Use(mw.Metrics)

	serverPort := ":" + os.Getenv("AUTHMS_HTTP_SERVER_PORT")
	logs.Logger.Info("starting auth http server at ", serverPort)

	err = http.ListenAndServe(serverPort, mainRouter)
	if err != nil {
		logs.LogFatal(logs.Logger, "auth", "main", err, err.Error())
	}
	logs.Logger.Info("auth http server stopped")
}
