package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/joho/godotenv"

	_http "2023_2_Holi/auth/delivery/http"
	"2023_2_Holi/auth/delivery/http/middleware"
	"2023_2_Holi/auth/repository/postgresql"
	"2023_2_Holi/auth/usecase"

	_ "github.com/lib/pq"
)

func fromEnv() string {
	host := os.Getenv("DB_HOST")
	if host == "" {
		return ""
	}

	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	dbname := os.Getenv("DB_NAME")

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, dbname)
}

func loggerInit() *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{})

	logLevel := os.Getenv("LOG_LEVEL")
	switch logLevel {
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
	case "info":
		logger.SetLevel(logrus.InfoLevel)
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
	case "fatal":
		logger.SetLevel(logrus.FatalLevel)
	}
	return logger
}

type AccessLogger struct {
	LogrusLogger *logrus.Entry
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
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Failed to get config : ", err)
	}

	logger := loggerInit()
	accessLogger := AccessLogger{
		LogrusLogger: logrus.NewEntry(logger),
	}

	logger.Info("starting connect to db")

	db, err := sql.Open("postgres", fromEnv())
	logger.Debug("db conf :", db)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"package":  "main",
			"function": "main",
			"error":    err,
		}).Fatal("Failed to open db")
	}
	defer db.Close()

	router := mux.NewRouter()
	router.Use(accessLogger.accessLogMiddleware)

	sessionRepository := postgresql.NewSessionPostgresqlRepository(db)
	authRepository := postgresql.NewAuthPostgresqlRepository(db)
	authUsecase := usecase.NewAuthUsecase(authRepository, sessionRepository)

	authRouter := mux.NewRouter()
	_http.NewAuthHandler(authRouter, authUsecase)

	mainRouter := authRouter.PathPrefix("api/").Subrouter()
	mw := middleware.InitMiddleware(authUsecase)
	mainRouter.Use(mw.IsAuth)

	logger.Info("starting server at :8080")
	err = http.ListenAndServe(":8080", authRouter)

	if err != nil {
		logger.WithFields(logrus.Fields{
			"package":  "main",
			"function": "main",
			"error":    err,
		}).Fatal("Failed to open server")
	}
	logger.Info("server stopped")
}
