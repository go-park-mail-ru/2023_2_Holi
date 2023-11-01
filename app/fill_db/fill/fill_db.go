package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

const pathPreview = "https://static_holi.hb.ru-msk.vkcs.cloud/Preview_Film/"
const pathMedia = "https://static_holi.hb.ru-msk.vkcs.cloud/Media_Files/"
const pathPreviewMedia = "https://static_holi.hb.ru-msk.vkcs.cloud/Media_Preview/"

func ageRes(age string) int {
	switch age {
	case "R":
		return 18
	case "PG-13":
		return 13
	case "PG":
		return 10
	case "TV-MA":
		return 18
	case "TV-14":
		return 14
	case "TV-7":
		return 7
	default:
		return 16
	}
}

func dbParamsfromEnv() string {
	host := os.Getenv("POSTGRES_HOST")
	host = "localhost"
	port := os.Getenv("POSTGRES_PORT")
	port = "5432"
	user := os.Getenv("POSTGRES_USER")
	user = "postgres"
	pass := os.Getenv("POSTGRES_PASSWORD")
	pass = "1784"
	dbname := os.Getenv("POSTGRES_DB")
	dbname = "netflix"

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

	file, err := os.Open("Netflix_Dataset.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	genreMap := make(map[string]int)

	castMap := make(map[string]int)

	castCount := 1

	genreCount := 1

	_, err = reader.Read()
	if err != nil {
		log.Fatal(err)
	}

	for {
		genreID++
		if genreID == 22 {
			break
		}
		row, err := reader.Read()
		if err != nil {
			break
		}
		genres := strings.Split(row[1], ",")
		for _, genre := range genres {
			genre = strings.TrimSpace(genre)

			idToInsert, genreExists := genreMap[genre]

			if !genreExists {
				er := db.QueryRow(`SELECT id FROM genre WHERE name = $1`, genre).Scan(&idToInsert)
				if er != nil {
					_, er = db.Exec(`INSERT INTO genre (id,name) VALUES ($1,$2)`, genreCount, genre)
					genreCount++
					if er != nil {
						continue
					}

					er = db.QueryRow(`SELECT id FROM genre WHERE name = $1`, genre).Scan(&idToInsert)
					if er != nil {
						continue
					}
				}
				genreMap[genre] = idToInsert
			}
		}
		casts := strings.Split(row[10], ",")
		for _, cast := range casts {
			cast = strings.TrimSpace(cast)

			idToInsert, castExists := castMap[cast]

			if !castExists {
				er := db.QueryRow(`SELECT id FROM "cast" WHERE name = $1`, cast).Scan(&idToInsert)
				if er != nil {
					_, er = db.Exec(`INSERT INTO "cast" (id, name) VALUES ($1, $2)`, castCount, cast)
					castCount++
					if er != nil {
						continue
					}

					er = db.QueryRow(`SELECT id FROM "cast" WHERE name = $1`, cast).Scan(&idToInsert)
					if er != nil {
						continue
					}
				}
				castMap[cast] = idToInsert
			}
		}
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		log.Fatal(err)
	}
	_, err = reader.Read()
	if err != nil {
		log.Fatal(err)
	}
	_, err = reader.Read()
	if err != nil {
		log.Fatal(err)
	}

	for {
		count++
		if count == 22 {
			break
		}
		row, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				count--
				continue
			}
		}
		genres := strings.Split(row[1], ",")
		casts := strings.Split(row[10], ",")

		sqlVideo := "INSERT INTO video (id, name, description, preview_video_path, release_year, rating, age_restriction, seasons_count) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"
		name := strings.Replace(row[0], " ", "_", -1)
		pr_pathMedia := pathPreviewMedia + name + ".mp4"
		release := row[19][:4]
		releaseInt, err := strconv.Atoi(release)
		if err != nil {
			fmt.Println("Ошибка при преобразовании в int:", err)
			return
		}
		age_restriction := ageRes(row[11])
		_, err = db.Exec(sqlVideo, count, row[0], row[23], pr_pathMedia, releaseInt, row[12], age_restriction, 1)
		if err != nil {
			log.Printf("Ошибка при вставке video: %v", err)
			continue
		}

		sqlEpisode := "INSERT INTO episode (id, name, description, duration ,preview_path, media_path, number, season_number, video_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)"
		name = strings.Replace(row[0], " ", "_", -1)
		pr_Path := pathPreview + name + ".jpg"
		pr_Media := pathMedia + name + ".mp4"
		duration := 0
		_, err = db.Exec(sqlEpisode, count, row[0], row[23], duration, pr_Path, pr_Media, 1, 1, count)
		if err != nil {
			log.Printf("Ошибка при вставке video: %v", err)
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

			sqlStatement := "INSERT INTO video_genre (genre_id, video_id) VALUES ($1, $2)"
			_, err = db.Exec(sqlStatement, genreID, count)
			if err != nil {
				log.Printf("Ошибка при вставке записи в genre_video: %v", err)
			}
		}

		for _, cast := range casts {
			cast = strings.TrimSpace(cast)
			var castID int
			err := db.QueryRow(`SELECT id FROM "cast" WHERE name = $1`, cast).Scan(&castID)
			if err != nil {
				log.Printf("Ошибка при получении ID касоа %s: %v", cast, err)
				continue
			}

			sqlStatement := "INSERT INTO video_cast (cast_id, video_id) VALUES ($1, $2)"
			_, err = db.Exec(sqlStatement, castID, count)
			if err != nil {
				log.Printf("Ошибка при вставке записи в cast_video: %v", err)
			}
		}
	}

	fmt.Println("Данные из CSV-файла успешно вставлены в таблицу films.")
}
