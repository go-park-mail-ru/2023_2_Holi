package main

import (
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	_http "2023_2_Holi/auth/delivery/http"
	"2023_2_Holi/auth/delivery/http/middleware"
	"2023_2_Holi/auth/repository/postgresql"
	session "2023_2_Holi/auth/repository/redis"
	"2023_2_Holi/auth/usecase"
	postgres "2023_2_Holi/db_connector/postgres"
	redis "2023_2_Holi/db_connector/redis"

	"2023_2_Holi/collections/collections_usecase"
	_httpCol "2023_2_Holi/collections/delivery/collections_http"
	"2023_2_Holi/collections/repository/collections_postgresql"
	logs "2023_2_Holi/logger"

	_ "github.com/lib/pq"
)

type AccessLogger struct {
	LogrusLogger *logrus.Logger
}

func (ac *AccessLogger) accessLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)

		ac.LogrusLogger.WithFields(logrus.Fields{
			"method":      r.Method,
			"remote_addr": r.RemoteAddr,
			"work_time":   time.Since(start),
		}).Info(r.URL.Path)
	})
}

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
	accessLogger := AccessLogger{
		LogrusLogger: logs.Logger,
	}

	postgres := postgres.PostgresConnector()
	defer postgres.Close()

	redis := redis.RedisConnector()
	defer redis.Close()

	mainRouter := mux.NewRouter()
	authMiddlewareRouter := mainRouter.PathPrefix("/api").Subrouter()

	sessionRepository := session.NewSessionRedisRepository(redis)
	authRepository := postgresql.NewAuthPostgresqlRepository(postgres)
	filmRepository := collections_postgresql.NewFilmPostgresqlRepository(postgres)

	authUsecase := usecase.NewAuthUsecase(authRepository, sessionRepository)
	filmUsecase := collections_usecase.NewFilmUsecase(filmRepository)

	_http.NewAuthHandler(authMiddlewareRouter, mainRouter, authUsecase)
	_httpCol.NewFilmHandler(authMiddlewareRouter, filmUsecase)

	mw := middleware.InitMiddleware(authUsecase)

	authMiddlewareRouter.Use(mw.IsAuth)
	mainRouter.Use(accessLogger.accessLogMiddleware)
	mainRouter.Use(mux.CORSMethodMiddleware(mainRouter))
	mainRouter.Use(mw.CORS)

	serverPort := ":" + os.Getenv("SERVER_PORT")
	logs.Logger.Info("starting server at", serverPort)

	err := http.ListenAndServe(":8080", mainRouter)
	if err != nil {
		logs.LogFatal(logs.Logger, "main", "main", err, "Failed to start server")
	}
	logs.Logger.Info("server stopped")
}
