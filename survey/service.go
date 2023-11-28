package survey

import (
	g_sess "2023_2_Holi/domain/grpc/session"
	"context"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"github.com/joho/godotenv"

	grpc_connector "2023_2_Holi/connectors/grpc"
	"2023_2_Holi/connectors/postgres"
	logs "2023_2_Holi/logger"
	"2023_2_Holi/middleware"

	_ "github.com/lib/pq"

	survey_http "2023_2_Holi/survey/delivery/http"
	survey_postgres "2023_2_Holi/survey/repository/postgresql"
	survey_usecase "2023_2_Holi/survey/usecase"
)

func StartServer() {
	err := godotenv.Load()
	ctx := context.Background()
	accessLogger := middleware.AccessLogger{
		LogrusLogger: logs.Logger,
	}

	dbParams := postgres.GetParamsForSurveyDB()
	pc := postgres.Connect(ctx, dbParams)
	defer pc.Close()

	//tokens, _ := domain.NewHMACHashToken("Gvjhlk123bl1lma0")

	mainRouter := mux.NewRouter()
	authMiddlewareRouter := mainRouter.PathPrefix("/api").Subrouter()

	sr := survey_postgres.NewSurveyPostgresqlRepository(pc, ctx)
	su := survey_usecase.NewSurveyUsecase(sr)
	survey_http.NewSurveyHandler(authMiddlewareRouter, su)
	//sanitizer := bluemonday.UGCPolicy()

	gc := grpc_connector.Connect(os.Getenv("AUTHMS_GRPC_SERVER_HOST") + ":" + os.Getenv("AUTHMS_GRPC_SERVER_PORT"))
	mw := middleware.InitMiddleware(g_sess.NewAuthCheckerClient(gc), nil)

	authMiddlewareRouter.Use(mw.IsAuth)
	mainRouter.Use(accessLogger.AccessLogMiddleware)
	mainRouter.Use(mux.CORSMethodMiddleware(mainRouter))
	mainRouter.Use(mw.CORS)
	// mainRouter.Use(mw.CSRFProtection)

	serverPort := ":" + os.Getenv("SURVAYMS_HTTP_SERVER_PORT")
	logs.Logger.Info("starting server at ", serverPort)

	err = http.ListenAndServe(serverPort, mainRouter)
	if err != nil {
		logs.LogFatal(logs.Logger, "main", "main", err, err.Error())
	}
	logs.Logger.Info("server stopped")
}
