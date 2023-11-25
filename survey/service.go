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
)

func StartServer() {
	err := godotenv.Load()
	ctx := context.Background()
	accessLogger := middleware.AccessLogger{
		LogrusLogger: logs.Logger,
	}

	pc := postgres.Connect(ctx)
	defer pc.Close()

	//rc := redis.Connect()
	//defer rc.Close()

	//tokens, _ := domain.NewHMACHashToken("Gvjhlk123bl1lma0")

	mainRouter := mux.NewRouter()
	authMiddlewareRouter := mainRouter.PathPrefix("/api").Subrouter()

	//srr := search_postgres.NewSearchPostgresqlRepository(pc, ctx)
	//sr := auth_redis.NewSessionRedisRepository(rc)
	//ur := utils_redis.NewUtilsRedisRepository(rc)
	//ar := auth_postgres.NewAuthPostgresqlRepository(pc, ctx)
	// fr := films_postgres.NewFilmsPostgresqlRepository(pc, ctx)
	// gr := genre_postgres.GenrePostgresqlRepository(pc, ctx)
	// pr := profile_postgres.NewProfilePostgresqlRepository(pc, ctx)
	//fvr := favourites_postgres.NewFavouritesPostgresqlRepository(pc, ctx)

	//au := auth_usecase.NewAuthUsecase(ar, sr)
	// fu := films_usecase.NewFilmsUsecase(fr)
	// gu := genre_usecase.NewGenreUsecase(gr)
	//uu := utils_usecase.NewUtilsUsecase(ur)
	//fvu := favourites_usecase.NewFavouritesUsecase(fvr)
	//su := search_usecase.NewSearchUsecase(srr)

	// sess, _ := session.NewSession()
	// svc := s3.New(sess, aws.NewConfig().WithEndpoint(vkCloudHotboxEndpoint).WithRegion(defaultRegion))
	// pu := profile_usecase.NewProfileUsecase(pr, svc)

	//sanitizer := bluemonday.UGCPolicy()

	gc := grpc_connector.Connect(os.Getenv("AUTHMS_GRPC_SERVER_HOST") + ":" + os.Getenv("AUTHMS_GRPC_SERVER_PORT"))
	mw := middleware.InitMiddleware(g_sess.NewAuthCheckerClient(gc), nil)

	authMiddlewareRouter.Use(mw.IsAuth)
	mainRouter.Use(accessLogger.AccessLogMiddleware)
	mainRouter.Use(mux.CORSMethodMiddleware(mainRouter))
	mainRouter.Use(mw.CORS)
	mainRouter.Use(mw.CSRFProtection)

	serverPort := ":" + os.Getenv("SURVAYMS_HTTP_SERVER_PORT")
	logs.Logger.Info("starting server at ", serverPort)

	err = http.ListenAndServe(serverPort, mainRouter)
	if err != nil {
		logs.LogFatal(logs.Logger, "main", "main", err, err.Error())
	}
	logs.Logger.Info("server stopped")
}
