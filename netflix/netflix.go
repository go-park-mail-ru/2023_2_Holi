package netflix

import (
	"context"
	"net/http"
	"os"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	auth_http "2023_2_Holi/auth/delivery/http"
	auth_postgres "2023_2_Holi/auth/repository/postgresql"
	auth_redis "2023_2_Holi/auth/repository/redis"
	auth_usecase "2023_2_Holi/auth/usecase"

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

	rc := redis.Connect()
	defer rc.Close()

	mainRouter := mux.NewRouter()
	authMiddlewareRouter := mainRouter.PathPrefix("/api").Subrouter()

	sessionRepository := auth_redis.NewSessionRedisRepository(rc)
	authRepository := auth_postgres.NewAuthPostgresqlRepository(pc, ctx)
	filmRepository := films_postgres.NewFilmsPostgresqlRepository(pc, ctx)
	genreRepository := genre_postgres.GenrePostgresqlRepository(pc, ctx)
	profileRepository := profile_postgres.NewProfilePostgresqlRepository(pc, ctx)

	authUsecase := auth_usecase.NewAuthUsecase(authRepository, sessionRepository)
	filmsUsecase := films_usecase.NewFilmsUsecase(filmRepository)
	genreUsecase := genre_usecase.NewGenreUsecase(genreRepository)
	profileUsecase := profile_usecase.NewProfileUsecase(profileRepository)

	auth_http.NewAuthHandler(authMiddlewareRouter, mainRouter, authUsecase)
	films_http.NewFilmsHandler(authMiddlewareRouter, filmsUsecase)
	genre_http.NewGenreHandler(authMiddlewareRouter, genreUsecase)
	profile_http.NewProfileHandler(authMiddlewareRouter, profileUsecase)

	mw := middleware.InitMiddleware(authUsecase)

	authMiddlewareRouter.Use(mw.IsAuth)
	mainRouter.Use(accessLogger.AccessLogMiddleware)
	mainRouter.Use(mux.CORSMethodMiddleware(mainRouter))
	mainRouter.Use(mw.CORS)
	mainRouter.Use(csrf.Protect([]byte("32-byte-long-auth-key"), csrf.Secure(false), csrf.HttpOnly(false), csrf.Path("/")))

	serverPort := ":" + os.Getenv("SERVER_PORT")
	logs.Logger.Info("starting server at ", serverPort)

	err = http.ListenAndServe(serverPort, mainRouter)
	if err != nil {
		logs.LogFatal(logs.Logger, "main", "main", err, err.Error())
	}
	logs.Logger.Info("server stopped")
}
