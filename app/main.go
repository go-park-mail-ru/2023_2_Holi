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

	"2023_2_Holi/collections/collections_usecase"
	_httpCol "2023_2_Holi/collections/delivery/collections_http"
	"2023_2_Holi/collections/repository/collections_postgresql"
	logs "2023_2_Holi/logs"

	_httpGen "2023_2_Holi/genre/delivery/genre_http"
	"2023_2_Holi/genre/genre_usecase"
	"2023_2_Holi/genre/repository/genre_postgresql"

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
	dbname := os.Getenv("POSTGRES_NAME")

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
// @schemes Zhttp
// @BasePath /
func main() {
	err := godotenv.Load()
	if err != nil {
		logs.Logger.Fatal("Failed to get config : ", err)
	}

	/*accessLogger := AccessLogger{
		LogrusLogger: logs.Logger,
	}*/

	db, err := sql.Open("postgres", dbParamsfromEnv())
	if err != nil {
		logs.LogFatal(logs.Logger, "main", "main", err, "Failed to connect to db")
	}
	defer db.Close()
	logs.Logger.Debug("db conf :", db)

	err = db.Ping()
	if err != nil {
		logs.LogFatal(logs.Logger, "main", "main", err, "DB doesn't listen")
	}
	logs.Logger.Info("Connected to postgres")

	mainRouter := mux.NewRouter()

	sessionRepository := postgresql.NewSessionPostgresqlRepository(db)
	authRepository := postgresql.NewAuthPostgresqlRepository(db)
	authUsecase := usecase.NewAuthUsecase(authRepository, sessionRepository)

	authMiddlewareRouter := mainRouter.PathPrefix("/api").Subrouter()

	_http.NewAuthHandler(authMiddlewareRouter, mainRouter, authUsecase)

	genreRepository := genre_postgresql.GenrePostgresqlRepository(db)
	genreUsecase := genre_usecase.NewGenreUsecase(genreRepository)
	_httpGen.NewGenreHandler(authMiddlewareRouter, genreUsecase)

	filmRepository := collections_postgresql.NewFilmPostgresqlRepository(db)
	filmUsecase := collections_usecase.NewFilmUsecase(filmRepository)
	_httpCol.NewFilmHandler(authMiddlewareRouter, filmUsecase)

	/*mw := middleware.InitMiddleware(authUsecase)

	authMiddlewareRouter.Use(mw.IsAuth)
	mainRouter.Use(accessLogger.accessLogMiddleware)
	mainRouter.Use(mux.CORSMethodMiddleware(mainRouter))
	mainRouter.Use(mw.CORS)*/

	logs.Logger.Info("starting server at :8080")

	err = http.ListenAndServe(":8080", mainRouter)
	if err != nil {
		logs.LogFatal(logs.Logger, "main", "main", err, "Failed to start server")
	}
	logs.Logger.Info("server stopped")
}
