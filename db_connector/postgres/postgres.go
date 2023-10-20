package postgres_connector

import (
	logs "2023_2_Holi/logger"
	"database/sql"
	"fmt"
	"os"
)

func PostgresConnector() *sql.DB {
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	pass := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")
	params := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, dbname)

	postgres, err := sql.Open("postgres", params)
	if err != nil {
		logs.LogFatal(logs.Logger, "postgres_connector", "PostgresConnector", err, "Failed to connect to postgres")
	}

	err = postgres.Ping()
	if err != nil {
		logs.LogFatal(logs.Logger, "postgres_connector", "PostgresConnector", err, "postgres doesn't listen")
	}
	logs.Logger.Info("Connected to postgres")
	logs.Logger.Debug("postgres client :", postgres)

	return postgres
}
