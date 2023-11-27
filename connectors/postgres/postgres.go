package postgres

import (
	"context"
	"fmt"
	"os"

	logs "2023_2_Holi/logger"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetParamsForNetflixDB() string {
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	pass := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")
	params := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, dbname)

	return params
}

//func GetParamsForSurveyDB() string {
//	host := os.Getenv("POSTGRES_CSAT_HOST")
//	port := "5432"
//	user := os.Getenv("POSTGRES_CSAT_USER")
//	pass := os.Getenv("POSTGRES_CSAT_PASSWORD")
//	dbname := os.Getenv("POSTGRES_CSAT_DB")
//	params := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, dbname)
//
//	return params
//}

func Connect(ctx context.Context, params string) *pgxpool.Pool {

	dbpool, err := pgxpool.New(ctx, params)
	if err != nil {
		logs.LogFatal(logs.Logger, "postgres_connector", "Connect", err, err.Error())
	}

	err = dbpool.Ping(ctx)
	if err != nil {
		logs.LogFatal(logs.Logger, "postgres_connector", "Connect", err, err.Error())
	}
	logs.Logger.Info("Connected to postgres")
	logs.Logger.Debug("postgres client :", dbpool)

	return dbpool
}
