package profile

import (
	grpc_connector "2023_2_Holi/connectors/grpc"
	"2023_2_Holi/connectors/postgres"
	g_sess "2023_2_Holi/domain/grpc/session"
	logs "2023_2_Holi/logger"
	"2023_2_Holi/middleware"
	profile_http "2023_2_Holi/profile/delivery/http"
	profile_postgres "2023_2_Holi/profile/repository/postgresql"
	profile_usecase "2023_2_Holi/profile/usecase"
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/microcosm-cc/bluemonday"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"os"
)

const (
	vkCloudHotboxEndpoint = "https://hb.vkcs.cloud"
	defaultRegion         = "ru-msk"
)

func StartService() {
	err := godotenv.Load()
	ctx := context.Background()
	accessLogger := middleware.AccessLogger{
		LogrusLogger: logs.Logger,
	}

	dbParams := postgres.GetParamsForNetflixDB()
	pc := postgres.Connect(ctx, dbParams)
	defer pc.Close()

	mainRouter := mux.NewRouter()
	authMiddlewareRouter := mainRouter.PathPrefix("/api").Subrouter()

	pr := profile_postgres.NewProfilePostgresqlRepository(pc, ctx)
	sess, _ := session.NewSession()
	svc := s3.New(sess, aws.NewConfig().WithEndpoint(vkCloudHotboxEndpoint).WithRegion(defaultRegion))
	pu := profile_usecase.NewProfileUsecase(pr, svc)
	sanitizer := bluemonday.UGCPolicy()

	profile_http.NewProfileHandler(authMiddlewareRouter, pu, sanitizer)
	mainRouter.Handle("/metrics", promhttp.Handler())

	gc := grpc_connector.Connect(os.Getenv("AUTHMS_GRPC_SERVER_HOST") + ":" + os.Getenv("AUTHMS_GRPC_SERVER_PORT"))
	mw := middleware.InitMiddleware(g_sess.NewAuthCheckerClient(gc), nil)

	mainRouter.Use(mw.Metrics)
	authMiddlewareRouter.Use(mw.IsAuth)
	mainRouter.Use(accessLogger.AccessLogMiddleware)
	mainRouter.Use(mux.CORSMethodMiddleware(mainRouter))
	mainRouter.Use(mw.CORS)
	mainRouter.Use(mw.CSRFProtection)

	serverPort := ":" + os.Getenv("PROFILEMS_HTTP_SERVER_PORT")
	logs.Logger.Info("starting service at ", serverPort)

	err = http.ListenAndServe(serverPort, mainRouter)
	if err != nil {
		logs.LogFatal(logs.Logger, "profile", "StartService", err, err.Error())
	}
	logs.Logger.Info("server stopped")
}
