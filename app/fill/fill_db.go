package fill

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

func main() {
	connStr := "host=localhost port=5432 dbname=netflix_proj user=aleksej password=postgres sslmode=disable"
	count := 0
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer db.Close()

	file, err := os.Open("app/fill/Netflix_Dataset.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	headers, err := reader.Read()
	if err != nil {
		log.Fatal(err)
	}
	for {
		count++
		if count == 101 {
			break
		}
		row, err := reader.Read()
		if err != nil {
			break
		}
		sqlStatement := "INSERT INTO film (id, name, preview_path,rating) VALUES ($1, $2, $3, $4)"
		strings.Join(headers, ", ")
		_, err = db.Exec(sqlStatement, count, row[0], row[26], row[12])
		if err != nil {
			log.Println(err)
			continue
		}
		genres := strings.Split(row[1], ",")

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
