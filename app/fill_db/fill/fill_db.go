package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

const path = "https://static_holi.hb.ru-msk.vkcs.cloud/Preview_Film/"

func dbParamsfromEnv() string {
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	pass := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, dbname)
}

func main() {
	count := 0
	genreID := 0
	db, err := sql.Open("postgres", dbParamsfromEnv())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}

	file, err := os.Open("../Netflix_Dataset.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	headers, err := reader.Read()
	if err != nil {
		log.Fatal(err)
	}

	genreMap := make(map[string]int)

	for {
		genreID++
		if genreID == 2000 {
			break
		}
		row, err := reader.Read()
		if err != nil {
			break
		}
		genres := strings.Split(row[1], ",")
		for _, genre := range genres {
			genre = strings.TrimSpace(genre)
			idToInsert := genreID
			if existingID, ok := genreMap[genre]; ok {
				idToInsert = existingID
			} else {
				er := db.QueryRow("SELECT id FROM genre WHERE name = $1", genre).Scan(&idToInsert)
				if er != nil {
					strings.Join(headers, ", ")
					_, er = db.Exec("INSERT INTO genre (id, name) VALUES ($1,$2)", idToInsert, genre)
					if er != nil {
						idToInsert--
						continue
					}
				}
				genreMap[genre] = idToInsert
			}
		}
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		log.Fatal(err)
	}

	for {
		count++
		if count == 100 {
			fmt.Println("break")
			break
		}
		row, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				continue
			}
		}
		genres := strings.Split(row[1], ",")

		sqlStatement := "INSERT INTO film (id, name, preview_path, rating) VALUES ($1, $2, $3, $4)"
		name := strings.Replace(row[0], " ", "_", -1)
		preview_path := path + name + ".jpg"
		_, err = db.Exec(sqlStatement, count, row[0], preview_path, row[12])
		if err != nil {
			log.Printf("Ошибка при вставке фильма: %v", err)
			continue
		}

		for _, genre := range genres {
			genre = strings.TrimSpace(genre)
			var genreID int
			err := db.QueryRow("SELECT id FROM genre WHERE name = $1", genre).Scan(&genreID)
			if err != nil {
				log.Printf("Ошибка при получении ID жанра %s: %v", genre, err)
				continue
			}

			sqlStatement := "INSERT INTO genre_film (genre_id, film_id) VALUES ($1, $2)"
			_, err = db.Exec(sqlStatement, genreID, count)
			if err != nil {
				log.Printf("Ошибка при вставке записи в genre-film: %v", err)
			}
		}
	}

	fmt.Println("Данные из CSV-файла успешно вставлены в таблицу films.")
}
