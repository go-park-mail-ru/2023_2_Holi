package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"

	_http "2023_2_Holi/auth/delivery/http"
	"2023_2_Holi/auth/delivery/http/middleware"
	"2023_2_Holi/auth/repository/postgresql"
	session "2023_2_Holi/auth/repository/redis"
	"2023_2_Holi/auth/usecase"

	"2023_2_Holi/collections/collections_usecase"
	_httpCol "2023_2_Holi/collections/delivery/collections_http"
	"2023_2_Holi/collections/repository/collections_postgresql"
	logs "2023_2_Holi/logs"

	_ "github.com/lib/pq"
)

func redisConnector() *redis.Client {
	redis := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
	})
	defer redis.Close()

	_, err := redis.Ping(context.Background()).Result()
	if err != nil {
		logs.LogFatal(logs.Logger, "main", "redisConnector", err, "Failed to connect to redis")
	}
	logs.Logger.Info("Connected to redis")

	return redis
}

func postgresParamsfromEnv() string {
	host := os.Getenv("POSTGRES_HOST")
	if host == "" {
		return ""
	}

	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	pass := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, dbname)
}

type AccessLogger struct {
	LogrusLogger *logrus.Logger
}

func (ac *AccessLogger) accessLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		//w.Header().Set("Access-Control-Allow-Origin", "*")

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

	postgres, err := sql.Open("postgres", postgresParamsfromEnv())
	if err != nil {
		logs.LogFatal(logs.Logger, "main", "main", err, "Failed to connect to postgres")
	}
	defer postgres.Close()
	logs.Logger.Debug("postgres conf :", postgres)

	redis := redisConnector()
	logs.Logger.Debug("redis client :", redis)

	err = postgres.Ping()
	if err != nil {
		logs.LogFatal(logs.Logger, "main", "main", err, "postgres doesn't listen")
	}
	logs.Logger.Info("Connected to postgres")

	mainRouter := mux.NewRouter()

	sessionRepository := session.NewSessionRedisRepository(redis)

	authRepository := postgresql.NewAuthPostgresqlRepository(postgres)
	authUsecase := usecase.NewAuthUsecase(authRepository, sessionRepository)

	authMiddlewareRouter := mainRouter.PathPrefix("/api").Subrouter()

	_http.NewAuthHandler(authMiddlewareRouter, mainRouter, authUsecase)

	filmRepository := collections_postgresql.NewFilmPostgresqlRepository(postgres)
	filmUsecase := collections_usecase.NewFilmUsecase(filmRepository)
	_httpCol.NewFilmHandler(authMiddlewareRouter, filmUsecase)

	mw := middleware.InitMiddleware(authUsecase)

	authMiddlewareRouter.Use(mw.IsAuth)
	mainRouter.Use(accessLogger.accessLogMiddleware)
	mainRouter.Use(mux.CORSMethodMiddleware(mainRouter))
	mainRouter.Use(mw.CORS)

	logs.Logger.Info("starting server at :8080")

	err = http.ListenAndServe(":8080", mainRouter)
	if err != nil {
		logs.LogFatal(logs.Logger, "main", "main", err, "Failed to start server")
	}
	logs.Logger.Info("server stopped")
}
