package subscription

import (
	grpc_connector "2023_2_Holi/connectors/grpc"
	"2023_2_Holi/connectors/postgres"
	"2023_2_Holi/domain"
	g_sess "2023_2_Holi/domain/grpc/session"
	"2023_2_Holi/domain/grpc/subscription"
	logs "2023_2_Holi/logger"
	"2023_2_Holi/middleware"
	sub_grpc "2023_2_Holi/subscription/delivery/grpc"
	"context"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"os"

	subscription_http "2023_2_Holi/subscription/delivery/http"
	subscription_postgres "2023_2_Holi/subscription/repository/postgresql"
	subscription_usecase "2023_2_Holi/subscription/usecase"
)

func startRpcServer(su domain.SubsUsecase) {
	port := ":" + os.Getenv("SUBMS_GRPC_SERVER_PORT")
	lis, err := net.Listen("tcp", port)
	if err != nil {
		logs.Logger.Info(logs.Logger, "subscription", "startRpcServer", err, err.Error())
	}

	server := grpc.NewServer()

	subscription.RegisterSubCheckerServer(server, sub_grpc.NewSubHandler(su))

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

	dbParams := postgres.GetParamsForUsrDB()
	pc := postgres.Connect(ctx, dbParams)
	defer pc.Close()

	mainRouter := mux.NewRouter()
	authMiddlewareRouter := mainRouter.PathPrefix("/api").Subrouter()

	subr := subscription_postgres.NewSubsPostgresqlRepository(pc, ctx)
	subu := subscription_usecase.NewSubsUsecase(subr)

	go startRpcServer(subu)

	subscription_http.NotAuthSubHandler(mainRouter, subu)
	subscription_http.NewSubsHandler(authMiddlewareRouter, subu)

	//mainRouter.Handle("/metrics", promhttp.Handler())

	gc := grpc_connector.Connect(os.Getenv("AUTHMS_GRPC_SERVER_HOST") + ":" + os.Getenv("AUTHMS_GRPC_SERVER_PORT"))
	mw := middleware.InitMiddleware(g_sess.NewAuthCheckerClient(gc), nil)

	//mainRouter.Use(mw.Metrics)
	authMiddlewareRouter.Use(mw.IsAuth)
	mainRouter.Use(accessLogger.AccessLogMiddleware)
	mainRouter.Use(mux.CORSMethodMiddleware(mainRouter))
	mainRouter.Use(mw.CORS)
	//mainRouter.Use(mw.CSRFProtection)

	serverPort := ":" + os.Getenv("SUBMS_HTTP_SERVER_PORT")
	logs.Logger.Info("starting service at ", serverPort)

	err = http.ListenAndServe(serverPort, mainRouter)
	if err != nil {
		logs.LogFatal(logs.Logger, "sub", "StartService", err, err.Error())
	}
	logs.Logger.Info("server stopped")
}
