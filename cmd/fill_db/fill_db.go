package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const pathPreview = "https://eyescreen.hb.ru-msk.vkcs.cloud/Preview_Film/"
const pathMedia = "https://eyescreen.hb.ru-msk.vkcs.cloud/Media_Files/"
const pathPreviewMedia = "https://eyescreen.hb.ru-msk.vkcs.cloud/Media_Preview/"
const actorsImg = "https://eyescreen.hb.ru-msk.vkcs.cloud/Actors_Image/"
const episodeIMG = "https://eyescreen.hb.ru-msk.vkcs.cloud/Episode_IMG/"

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

func generateRandomRating() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(10) + 1
}

func dbParamsfromEnvUsr() string {
	host := os.Getenv("POSTGRES_USR_HOST")
	port := os.Getenv("POSTGRES_USR_PORT")
	user := os.Getenv("POSTGRES_USR_USER")
	pass := os.Getenv("POSTGRES_USR_PASSWORD")
	dbname := os.Getenv("POSTGRES_USR_DB")
	host = "postgres_usr"
	port = "5432"
	user = "postgres"
	pass = "123"
	dbname = "netflix_auth"

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, dbname)
}

func dbParamsfromEnv() string {
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	pass := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, dbname)
}

func main() {
	err := godotenv.Load()
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

	fileActors, err := os.Open("Actors.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer fileActors.Close()

	readerActors := csv.NewReader(fileActors)
	readerActors.Comma = ';'
	readerActors.LazyQuotes = true

	_, err = readerActors.Read()
	if err != nil {
		log.Fatal(err)
	}

	for actorID := 1; actorID <= 157; actorID++ {
		row, err := readerActors.Read()
		if err != nil {
			log.Fatal(err)
		}

		name := row[0]
		birthday := row[1]
		place := row[2]
		career := row[3]
		img := actorsImg + strconv.Itoa(actorID) + ".jpeg"

		_, err = db.Exec(`INSERT INTO "cast" (id, name, birthday, place, carier, imgPath) VALUES ($1, $2, $3, $4, $5, $6)`,
			actorID, name, birthday, place, career, img)
		if err != nil {
			log.Printf("Ошибка при вставке актера: %v", err)
			continue
		}
	}

	file, err := os.Open("Netflix_Dataset.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	genreMap := make(map[string]int)

	genreCount := 1

	_, err = reader.Read()
	if err != nil {
		log.Fatal(err)
	}

	for {
		genreID++
		if genreID == 60 {
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
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		log.Fatal(err)
	}

	row, err := reader.Read()
	if err != nil {
		fmt.Println(row)
	}

	for {
		count++
		if count == 38 {
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

		sqlVideo := "INSERT INTO video (id, name, description, preview_path ,preview_video_path, release_year, rating, age_restriction, seasons_count) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)"
		name := strings.Replace(row[5], " ", "_", -1)
		pr_Path := pathPreview + name + ".jpeg"
		pr_pathMedia := pathPreviewMedia + name + ".mp4"
		release := row[19][:4]
		releaseInt, err := strconv.Atoi(release)
		if err != nil {
			fmt.Println("Ошибка при преобразовании в int:", err)
			return
		}
		age_restriction := ageRes(row[11])
		_, err = db.Exec(sqlVideo, count, row[0], row[23], pr_Path, pr_pathMedia, releaseInt, row[12], age_restriction, 0)
		if err != nil {
			log.Printf("Ошибка при вставке video: %v", err)
			continue
		}

		sqlEpisode := "INSERT INTO episode (id, name, description, duration ,preview_path, media_path, number, season_number, video_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)"
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

	i := 0
	count--
	countEpisode := 39
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading CSV:", err)
		return
	}
	for {
		count++
		if count == 39 {
			break
		}
		genres := strings.Split(records[i][1], ",")
		casts := strings.Split(records[i][10], ",")

		sqlVideo := "INSERT INTO video (id, name, description, preview_path ,preview_video_path, release_year, rating, age_restriction, seasons_count) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)"
		name := strings.Replace(records[i][5], " ", "_", -1)
		pr_Path := pathPreview + name + ".jpeg"
		pr_pathMedia := pathPreviewMedia + name + ".mp4"
		fmt.Println(count)
		release := records[i][19][:4]
		releaseInt, err := strconv.Atoi(release)
		if err != nil {
			fmt.Println("Ошибка при преобразовании в int:", err)
			return
		}
		age_restriction := ageRes(records[i][11])
		_, err = db.Exec(sqlVideo, count, records[i][0], records[i][23], pr_Path, pr_pathMedia, releaseInt, records[i][12], age_restriction, records[i][2])
		if err != nil {
			log.Printf("Ошибка при вставке video: %v", err)
			continue
		}

		i++

		for {
			if records[i][4] != "Episode" {
				break
			}
			fmt.Println(records[i][3])
			sqlEpisode := "INSERT INTO episode (id, name, description, duration ,preview_path, media_path, number, season_number, video_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)"
			pr_Path := episodeIMG + name + "_S" + records[i][2] + "_E" + records[i][3] + ".jpeg"
			pr_Media := pathMedia + name + "_S" + records[i][2] + "_E" + records[i][3] + ".mp4"
			duration := 0
			_, err = db.Exec(sqlEpisode, countEpisode, records[i][0], records[i][23], duration, pr_Path, pr_Media, records[i][3], records[i][2], count)
			if err != nil {
				log.Printf("Ошибка при вставке video: %v", err)
				continue
			}
			i++
			countEpisode++
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

	_, err = db.Exec(`UPDATE "video" SET tsv = setweight(to_tsvector(name), 'A');`)
	if err != nil {
		log.Printf("Ошибка при создании вектора для video.name: %v", err)
	}
	_, err = db.Exec(`UPDATE "cast" SET tsv = setweight(to_tsvector(name), 'A');`)
	if err != nil {
		log.Printf("Ошибка при создании вектора для cast.name: %v", err)
	}

	dbUsr, err := sql.Open("postgres", dbParamsfromEnvUsr())
	if err != nil {
		log.Fatal(err)
	}
	defer dbUsr.Close()

	err = dbUsr.Ping()
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}

	for userID := 30; userID <= 40; userID++ {
		username := fmt.Sprintf("user%d", userID)
		email := fmt.Sprintf("user%d@example.com", userID)
		password := fmt.Sprintf("password%d", userID)

		_, err := dbUsr.Exec(`INSERT INTO "user" (id ,name, email, password) VALUES ($1, $2, $3, $4)`, userID, username, email, password)
		if err != nil {
			log.Printf("Ошибка при вставке пользователя: %v", err)
			continue
		}

		for videoID := 1; videoID <= 38; videoID++ {
			rating := generateRandomRating()
			_, err := db.Exec(`INSERT INTO video_estimation (rate, video_id, user_id) VALUES ($1, $2, $3)`, rating, videoID, userID)
			if err != nil {
				log.Printf("Ошибка при вставке оценки: %v", err)
				continue
			}
		}
	}

	fmt.Println("Данные из CSV-файла успешно вставлены в таблицу films.")
}
