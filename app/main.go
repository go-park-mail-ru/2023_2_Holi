package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/joho/godotenv"

	_http "2023_2_Holi/auth/delivery/http"
	"2023_2_Holi/auth/repository/postgresql"
	"2023_2_Holi/auth/usecase"
	"2023_2_Holi/logfuncs"

	_ "github.com/lib/pq"
)

func dbParamsfromEnv() string {
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

var logger = logfuncs.LoggerInit()

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
// @schemes Zhttp
// @BasePath /
func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		logger.Fatal("Failed to get config : ", err)
	}

	accessLogger := AccessLogger{
		LogrusLogger: logger,
	}

	db, err := sql.Open("postgres", dbParamsfromEnv())
	if err != nil {
		logfuncs.LogFatal(logger, "main", "main", err, "Failed to connect to db")
	}
	defer db.Close()
	logger.Debug("db conf :", db)

	err = db.Ping()
	if err != nil {
		logfuncs.LogFatal(logger, "main", "main", err, "DB doesn't listen")
	}
	logger.Info("Connected to postgres")

	router := mux.NewRouter()
	router.Use(accessLogger.accessLogMiddleware)
	sessionRepository := postgresql.NewSessionPostgresqlRepository(db)
	authRepository := postgresql.NewAuthPostgresqlRepository(db)
	authUsecase := usecase.NewAuthUsecase(authRepository, sessionRepository)

	_http.NewAuthHandler(router, authUsecase)

	logger.Info("starting server at :8080")
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		logfuncs.LogFatal(logger, "main", "main", err, "Failed to start server")
	}
	logger.Info("server stopped")
}
