package netflix

import (
	g_sess "2023_2_Holi/domain/grpc/session"
	favourites_http "2023_2_Holi/favourites/delivery/http"
	"context"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/microcosm-cc/bluemonday"

	films_http "2023_2_Holi/films/delivery/http"
	films_postgres "2023_2_Holi/films/repository/postgresql"
	films_usecase "2023_2_Holi/films/usecase"

	favourites_postgres "2023_2_Holi/favourites/repository/postgresql"
	favourites_usecase "2023_2_Holi/favourites/usecase"

	genre_http "2023_2_Holi/genre/delivery/http"
	genre_postgres "2023_2_Holi/genre/repository/postgresql"
	profile_postgres "2023_2_Holi/profile/repository/postgresql"

	genre_usecase "2023_2_Holi/genre/usecase"

	profile_http "2023_2_Holi/profile/delivery/http"
	profile_usecase "2023_2_Holi/profile/usecase"

	grpc_connector "2023_2_Holi/connectors/grpc"
	"2023_2_Holi/connectors/postgres"
	logs "2023_2_Holi/logger"
	"2023_2_Holi/middleware"
)

func StartServer() {
	err := godotenv.Load()
	ctx := context.Background()
	accessLogger := middleware.AccessLogger{
		LogrusLogger: logs.Logger,
	}

	dbParams := postgres.GetParamsForNetflixDB()
	pc := postgres.Connect(ctx, dbParams)
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
	fr := films_postgres.NewFilmsPostgresqlRepository(pc, ctx)
	gr := genre_postgres.GenrePostgresqlRepository(pc, ctx)
	pr := profile_postgres.NewProfilePostgresqlRepository(pc, ctx)
	fvr := favourites_postgres.NewFavouritesPostgresqlRepository(pc, ctx)

	//au := auth_usecase.NewAuthUsecase(ar, sr)
	fu := films_usecase.NewFilmsUsecase(fr)
	gu := genre_usecase.NewGenreUsecase(gr)
	//uu := utils_usecase.NewUtilsUsecase(ur)
	fvu := favourites_usecase.NewFavouritesUsecase(fvr)
	//su := search_usecase.NewSearchUsecase(srr)

	sess, _ := session.NewSession()
	svc := s3.New(sess, aws.NewConfig().WithEndpoint(vkCloudHotboxEndpoint).WithRegion(defaultRegion))
	pu := profile_usecase.NewProfileUsecase(pr, svc)

	sanitizer := bluemonday.UGCPolicy()

	//auth_http.NewAuthHandler(authMiddlewareRouter, mainRouter, au)
	films_http.NewFilmsHandler(authMiddlewareRouter, fu)
	genre_http.NewGenreHandler(authMiddlewareRouter, gu)
	profile_http.NewProfileHandler(authMiddlewareRouter, pu, sanitizer)
	//search_http.NewSearchHandler(authMiddlewareRouter, su)
	//csrf_http.NewCsrfHandler(mainRouter, tokens)
	//favourites_http.NewFavouritesHandler(authMiddlewareRouter, fvu, uu)

	gc := grpc_connector.Connect(os.Getenv("AUTHMS_GRPC_SERVER_HOST") + ":" + os.Getenv("AUTHMS_GRPC_SERVER_PORT"))
	mw := middleware.InitMiddleware(g_sess.NewAuthCheckerClient(gc), nil)

	authMiddlewareRouter.Use(mw.IsAuth)
	mainRouter.Use(accessLogger.AccessLogMiddleware)
	mainRouter.Use(mux.CORSMethodMiddleware(mainRouter))
	mainRouter.Use(mw.CORS)
	mainRouter.Use(mw.CSRFProtection)
	//mainRouter.Use(mw.Metrics)

	serverPort := ":" + os.Getenv("SERVER_PORT")
	logs.Logger.Info("starting server at ", serverPort)

	err = http.ListenAndServe(serverPort, mainRouter)
	if err != nil {
		logs.LogFatal(logs.Logger, "main", "main", err, err.Error())
	}
	logs.Logger.Info("server stopped")
}
