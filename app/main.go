package main

import (
	_http "2023_2_Holi/auth/delivery/http"
	"2023_2_Holi/auth/repository/postgresql"
	"2023_2_Holi/auth/usecase"
	"net/http"

	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
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

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("starting connect to db")
	db, err := sql.Open("postgres", fromEnv())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_http.NewAuthHandler(usecase.NewUserUsecase(postgresql.NewUserPostgresqlRepository(db)))

	fmt.Println("starting server at :8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("server stopped")
}
