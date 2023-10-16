package main

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"

	auth_http "2023_2_Holi/auth/delivery/http"
	auth_postgresql "2023_2_Holi/auth/repository/postgresql"
	auth_redis "2023_2_Holi/auth/repository/redis"
	auth_usecase "2023_2_Holi/auth/usecase"

	movies_http "2023_2_Holi/movies/delivery/http"
	movies_postgresql "2023_2_Holi/movies/repository/postgresql"
	movies_usecase "2023_2_Holi/movies/usecase"

	postgres "2023_2_Holi/db_connector/postgres"
	redis "2023_2_Holi/db_connector/redis"
	logs "2023_2_Holi/logger"
	middleware "2023_2_Holi/middleware"

	_ "github.com/lib/pq"
)

// @title Netfilx API
// @version 1.0
// @description API of the nelfix project by holi

// @contact.name Alex Chinaev
// @contact.url https://vk.com/l.chinaev
// @contact.email ax.chinaev@yandex.ru

// @license.name AS IS (NO WARRANTY)

// @host 127.0.0.1E
// @schemes http
// @BasePath /
func main() {
	accessLogger := middleware.AccessLogger{
		LogrusLogger: logs.Logger,
	}

	postgres := postgres.PostgresConnector()
	defer postgres.Close()

	redis := redis.RedisConnector()
	defer redis.Close()

	mainRouter := mux.NewRouter()
	authMiddlewareRouter := mainRouter.PathPrefix("/api").Subrouter()

	sessionRepository := auth_redis.NewSessionRedisRepository(redis)
	authRepository := auth_postgresql.NewAuthPostgresqlRepository(postgres)
	filmRepository := movies_postgresql.NewFilmPostgresqlRepository(postgres)

	authUsecase := auth_usecase.NewAuthUsecase(authRepository, sessionRepository)
	filmUsecase := movies_usecase.NewMoviesUsecase(filmRepository)

	auth_http.NewAuthHandler(authMiddlewareRouter, mainRouter, authUsecase)
	movies_http.NewMoviesHandler(authMiddlewareRouter, filmUsecase)

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
