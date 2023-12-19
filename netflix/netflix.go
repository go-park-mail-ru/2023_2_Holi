package netflix

import (
	"2023_2_Holi/connectors/redis"
	g_sess "2023_2_Holi/domain/grpc/session"
	//favourites_http "frontend/favourites/delivery/http"
	"context"
	"embed"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	genre_http "2023_2_Holi/genre/delivery/http"
	genre_postgres "2023_2_Holi/genre/repository/postgresql"
	genre_usecase "2023_2_Holi/genre/usecase"

	rating_http "2023_2_Holi/rating/delivery/http"
	rating_postgres "2023_2_Holi/rating/repository/postgresql"
	rating_usecase "2023_2_Holi/rating/usecase"

	search_http "2023_2_Holi/search/delivery/http"
	search_postgres "2023_2_Holi/search/repository/postgresql"
	search_usecase "2023_2_Holi/search/usecase"

	favourites_http "2023_2_Holi/favourites/delivery/http"
	favourites_postgres "2023_2_Holi/favourites/repository/postgresql"
	favourites_usecase "2023_2_Holi/favourites/usecase"

	grpc_connector "2023_2_Holi/connectors/grpc"
	"2023_2_Holi/connectors/postgres"
	logs "2023_2_Holi/logger"
	"2023_2_Holi/middleware"
)

const (
	vkCloudHotboxEndpoint = "https://hb.vkcs.cloud"
	defaultRegion         = "ru-msk"
)

var static embed.FS

func StartServer() {
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

	//tokens, _ := domain.NewHMACHashToken("Gvjhlk123bl1lma0")

	mainRouter := mux.NewRouter()
	authMiddlewareRouter := mainRouter.PathPrefix("/api").Subrouter()

	//srr := search_postgres.NewSearchPostgresqlRepository(pc, ctx)
	//sr := auth_redis.NewSessionRedisRepository(rc)
	//ur := utils_redis.NewUtilsRedisRepository(rc)
	//ar := auth_postgres.NewAuthPostgresqlRepository(pc, ctx)
	//fr := films_postgres.NewFilmsPostgresqlRepository(pc, ctx)
	gr := genre_postgres.GenrePostgresqlRepository(pc, ctx)
	rr := rating_postgres.NewRatingPostgresqlRepository(pc, ctx)
	sr := search_postgres.NewSearchPostgresqlRepository(pc, ctx)
	//pr := profile_postgres.NewProfilePostgresqlRepository(pc, ctx)
	fvr := favourites_postgres.NewFavouritesPostgresqlRepository(pc, ctx)
	//serr := series_postgres.NewSeriesPostgresqlRepository(pc, ctx)

	//au := auth_usecase.NewAuthUsecase(ar, sr)
	//fu := films_usecase.NewFilmsUsecase(fr)
	gu := genre_usecase.NewGenreUsecase(gr)
	ru := rating_usecase.NewRatingUsecase(rr)
	su := search_usecase.NewSearchUsecase(sr)
	//uu := utils_usecase.NewUtilsUsecase(ur)
	fvu := favourites_usecase.NewFavouritesUsecase(fvr)
	//seru := series_usecase.NewSeriesUsecase(serr)
	//su := search_usecase.NewSearchUsecase(srr)

	//sess, _ := session.NewSession()
	//svc := s3.New(sess, aws.NewConfig().WithEndpoint(vkCloudHotboxEndpoint).WithRegion(defaultRegion))
	//pu := profile_usecase.NewProfileUsecase(pr, svc)

	//sanitizer := bluemonday.UGCPolicy()

	//auth_http.NewAuthHandler(authMiddlewareRouter, mainRouter, au)
	//films_http.NewFilmsHandler(authMiddlewareRouter, fu)
	genre_http.NewGenreHandler(authMiddlewareRouter, gu)
	rating_http.NewRatingHandler(authMiddlewareRouter, ru)
	search_http.NewSearchHandler(authMiddlewareRouter, su)
	//profile_http.NewProfileHandler(authMiddlewareRouter, pu, sanitizer)
	//series_http.NewSeriesHandler(authMiddlewareRouter, seru)
	//search_http.NewSearchHandler(authMiddlewareRouter, su)
	//csrf_http.NewCsrfHandler(mainRouter, tokens)
	favourites_http.NewFavouritesHandler(authMiddlewareRouter, fvu)

	gc := grpc_connector.Connect(os.Getenv("AUTHMS_GRPC_SERVER_HOST") + ":" + os.Getenv("AUTHMS_GRPC_SERVER_PORT"))
	mw := middleware.InitMiddleware(g_sess.NewAuthCheckerClient(gc), nil)

	//authMiddlewareRouter.Use(mw.IsAuth)
	mainRouter.Use(accessLogger.AccessLogMiddleware)
	mainRouter.Use(mux.CORSMethodMiddleware(mainRouter))
	mainRouter.Use(mw.CORS)
	//mainRouter.Use(mw.CSRFProtection)
	//mainRouter.Use(mw.Metrics)

	serverPort := ":" + os.Getenv("SERVER_PORT")
	logs.Logger.Info("starting server at ", serverPort)

	err = http.ListenAndServe(serverPort, mainRouter)
	if err != nil {
		logs.LogFatal(logs.Logger, "main", "main", err, err.Error())
	}
	logs.Logger.Info("server stopped")
}
