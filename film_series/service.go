package film_series

import (
	grpc_connector "2023_2_Holi/connectors/grpc"
	"2023_2_Holi/connectors/postgres"
	g_sess "2023_2_Holi/domain/grpc/session"
	g_sub "2023_2_Holi/domain/grpc/subscription"
	films_http "2023_2_Holi/film_series/films/delivery/http"
	films_postgres "2023_2_Holi/film_series/films/repository/postgresql"
	films_usecase "2023_2_Holi/film_series/films/usecase"
	series_http "2023_2_Holi/film_series/series/delivery/http"
	series_postgres "2023_2_Holi/film_series/series/repository/postgresql"
	series_usecase "2023_2_Holi/film_series/series/usecase"
	logs "2023_2_Holi/logger"
	"2023_2_Holi/middleware"
	"context"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"os"
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

	fr := films_postgres.NewFilmsPostgresqlRepository(pc, ctx)
	fu := films_usecase.NewFilmsUsecase(fr)
	subCli := grpc_connector.Connect(os.Getenv("SUBMS_GRPC_SERVER_HOST") + ":" + os.Getenv("SUBMS_GRPC_SERVER_PORT"))
	films_http.NewFilmsHandler(authMiddlewareRouter, fu, g_sub.NewSubCheckerClient(subCli))

	serr := series_postgres.NewSeriesPostgresqlRepository(pc, ctx)
	seru := series_usecase.NewSeriesUsecase(serr)
	series_http.NewSeriesHandler(authMiddlewareRouter, seru)

	mainRouter.Handle("/metrics", promhttp.Handler())

	gc := grpc_connector.Connect(os.Getenv("AUTHMS_GRPC_SERVER_HOST") + ":" + os.Getenv("AUTHMS_GRPC_SERVER_PORT"))

	mw := middleware.InitMiddleware(g_sess.NewAuthCheckerClient(gc), nil)

	mainRouter.Use(mw.Metrics)
	authMiddlewareRouter.Use(mw.IsAuth)
	mainRouter.Use(accessLogger.AccessLogMiddleware)
	mainRouter.Use(mux.CORSMethodMiddleware(mainRouter))
	mainRouter.Use(mw.CORS)
	mainRouter.Use(mw.CSRFProtection)

	serverPort := ":" + os.Getenv("FILM_SERIES_HTTP_SERVER_PORT")
	logs.Logger.Info("starting service at ", serverPort)

	err = http.ListenAndServe(serverPort, mainRouter)
	if err != nil {
		logs.LogFatal(logs.Logger, "films_series", "StartService", err, err.Error())
	}
	logs.Logger.Info("server stopped")
}
