package main

import (
	"context"
	"net/http"
	"os"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"

	auth_http "2023_2_Holi/auth/delivery/http"
	auth_postgres "2023_2_Holi/auth/repository/postgresql"
	auth_redis "2023_2_Holi/auth/repository/redis"
	auth_usecase "2023_2_Holi/auth/usecase"

	films_http "2023_2_Holi/films/delivery/http"
	films_postgres "2023_2_Holi/films/repository/postgresql"
	films_usecase "2023_2_Holi/films/usecase"

	postgres "2023_2_Holi/db/connector/postgres"
	redis "2023_2_Holi/db/connector/redis"
	logs "2023_2_Holi/logger"
	middleware "2023_2_Holi/middleware"

	genre_http "2023_2_Holi/genre/delivery/http"
	genre_postgres "2023_2_Holi/genre/repository/postgresql"
	genre_usecase "2023_2_Holi/genre/usecase"

	_ "github.com/lib/pq"
)

// @title Netfilx API
// @version 1.0
// @description API of the nelfix project by holi

// @contact.name Alex Chinaev
// @contact.url https://vk.com/l.chinaev
// @contact.email ax.chinaev@yandex.ru

// @license.name AS IS (NO WARRANTY)

// @host 127.0.0.1
// @schemes http
// @BasePath /
func main() {
	ctx := context.Background()
	accessLogger := middleware.AccessLogger{
		LogrusLogger: logs.Logger,
	}

	postgres := postgres.PostgresConnector(ctx)
	defer postgres.Close()

	redis := redis.RedisConnector()
	defer redis.Close()

	csrfMiddleware := csrf.Protect([]byte("32-byte-long-auth-key"))
	mainRouter := mux.NewRouter()
	authMiddlewareRouter := mainRouter.PathPrefix("/api").Subrouter()
	authMiddlewareRouter.Use(csrfMiddleware)

	sessionRepository := auth_redis.NewSessionRedisRepository(redis)
	authRepository := auth_postgres.NewAuthPostgresqlRepository(postgres, ctx)
	filmRepository := films_postgres.NewFilmsPostgresqlRepository(postgres, ctx)
	genreRepository := genre_postgres.GenrePostgresqlRepository(postgres, ctx)

	authUsecase := auth_usecase.NewAuthUsecase(authRepository, sessionRepository)
	filmsUsecase := films_usecase.NewFilmsUsecase(filmRepository)
	genreUsecase := genre_usecase.NewGenreUsecase(genreRepository)

	auth_http.NewAuthHandler(authMiddlewareRouter, mainRouter, authUsecase)
	films_http.NewFilmsHandler(authMiddlewareRouter, filmsUsecase)
	genre_http.NewGenreHandler(authMiddlewareRouter, genreUsecase)

	mw := middleware.InitMiddleware(authUsecase)

	authMiddlewareRouter.Use(mw.IsAuth)
	mainRouter.Use(accessLogger.AccessLogMiddleware)
	mainRouter.Use(mux.CORSMethodMiddleware(mainRouter))
	mainRouter.Use(mw.CORS)

	serverPort := ":" + os.Getenv("SERVER_PORT")
	logs.Logger.Info("starting server at ", serverPort)

	err := http.ListenAndServe(serverPort, mainRouter)
	if err != nil {
		logs.LogFatal(logs.Logger, "main", "main", err, "Failed to start server")
	}
	logs.Logger.Info("server stopped")
}
