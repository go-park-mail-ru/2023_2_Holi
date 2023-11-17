package netflix

import (
	"context"
	"embed"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gorilla/mux"

	"github.com/joho/godotenv"
	"github.com/microcosm-cc/bluemonday"

	auth_http "2023_2_Holi/auth/delivery/http"
	auth_postgres "2023_2_Holi/auth/repository/postgresql"
	auth_redis "2023_2_Holi/auth/repository/redis"
	auth_usecase "2023_2_Holi/auth/usecase"
	"2023_2_Holi/domain"

	films_http "2023_2_Holi/films/delivery/http"
	films_postgres "2023_2_Holi/films/repository/postgresql"
	films_usecase "2023_2_Holi/films/usecase"

	"2023_2_Holi/db/connector/postgres"
	"2023_2_Holi/db/connector/redis"
	logs "2023_2_Holi/logger"
	"2023_2_Holi/middleware"

	genre_http "2023_2_Holi/genre/delivery/http"
	genre_postgres "2023_2_Holi/genre/repository/postgresql"
	genre_usecase "2023_2_Holi/genre/usecase"

	profile_http "2023_2_Holi/profile/delivery/http"
	profile_postgres "2023_2_Holi/profile/repository/postgresql"
	profile_usecase "2023_2_Holi/profile/usecase"

	search_http "2023_2_Holi/search/delivery/http"
	search_postgres "2023_2_Holi/search/repository/postgresql"
	search_usecase "2023_2_Holi/search/usecase"

	csrf_http "2023_2_Holi/csrf/delivery/http"

	_ "github.com/lib/pq"
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

	pc := postgres.Connect(ctx)
	defer pc.Close()

	rc := redis.Connect()
	defer rc.Close()

	tokens, _ := domain.NewHMACHashToken("Gvjhlk123bl1lma0")

	mainRouter := mux.NewRouter()
	authMiddlewareRouter := mainRouter.PathPrefix("/api").Subrouter()

	sessionRepository := auth_redis.NewSessionRedisRepository(rc)
	authRepository := auth_postgres.NewAuthPostgresqlRepository(pc, ctx)
	filmRepository := films_postgres.NewFilmsPostgresqlRepository(pc, ctx)
	genreRepository := genre_postgres.GenrePostgresqlRepository(pc, ctx)
	profileRepository := profile_postgres.NewProfilePostgresqlRepository(pc, ctx)
	searchRepository := search_postgres.NewSearchPostgresqlRepository(pc, ctx)

	authUsecase := auth_usecase.NewAuthUsecase(authRepository, sessionRepository)
	filmsUsecase := films_usecase.NewFilmsUsecase(filmRepository)
	genreUsecase := genre_usecase.NewGenreUsecase(genreRepository)
	searchUsecase := search_usecase.NewSearchUsecase(searchRepository)

	sess, _ := session.NewSession()
	svc := s3.New(sess, aws.NewConfig().WithEndpoint(vkCloudHotboxEndpoint).WithRegion(defaultRegion))
	profileUsecase := profile_usecase.NewProfileUsecase(profileRepository, svc)
	sanitizer := bluemonday.UGCPolicy()

	auth_http.NewAuthHandler(authMiddlewareRouter, mainRouter, authUsecase)
	films_http.NewFilmsHandler(authMiddlewareRouter, filmsUsecase)
	genre_http.NewGenreHandler(authMiddlewareRouter, genreUsecase)
	profile_http.NewProfileHandler(authMiddlewareRouter, profileUsecase, sanitizer)
	search_http.NewSearchHandler(authMiddlewareRouter, searchUsecase)
	csrf_http.NewCsrfHandler(mainRouter, tokens)

	mw := middleware.InitMiddleware(authUsecase)

	// authMiddlewareRouter.Use(mw.IsAuth)
	mainRouter.Use(accessLogger.AccessLogMiddleware)
	mainRouter.Use(mux.CORSMethodMiddleware(mainRouter))
	mainRouter.Use(mw.CORS)
	mainRouter.Use(mw.CSRFProtection)

	serverPort := ":" + os.Getenv("SERVER_PORT")
	logs.Logger.Info("starting server at ", serverPort)

	err = http.ListenAndServe(serverPort, mainRouter)
	if err != nil {
		logs.LogFatal(logs.Logger, "main", "main", err, err.Error())
	}
	logs.Logger.Info("server stopped")
}
