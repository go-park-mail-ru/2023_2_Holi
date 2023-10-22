package postgres_connector

import (
	"context"
	"fmt"
	"os"

	logs "2023_2_Holi/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

func PostgresConnector(ctx context.Context) *pgxpool.Pool {
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	pass := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")
	params := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, dbname)

	//postgres, err := sql.Open("postgres", params)
	//if err != nil {
	//	logs.LogFatal(logs.Logger, "postgres_connector", "PostgresConnector", err, "Failed to connect to postgres")
	//}

	dbpool, err := pgxpool.New(ctx, params)
	if err != nil {
		logs.LogFatal(logs.Logger, "postgres_connector", "PostgresConnector", err, "postgres doesn't listen")
	}
	defer dbpool.Close()

	err = dbpool.Ping(ctx)
	if err != nil {
		logs.LogFatal(logs.Logger, "postgres_connector", "PostgresConnector", err, "postgres doesn't listen")
	}
	logs.Logger.Info("Connected to postgres")
	logs.Logger.Debug("postgres client :", dbpool)

	return dbpool
}
