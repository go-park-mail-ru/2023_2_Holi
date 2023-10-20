package films_postgres

import (
	"database/sql"

	"2023_2_Holi/domain"
	logs "2023_2_Holi/logger"
)

type filmsPostgresqlRepository struct {
	db *sql.DB
}

func NewFilmsPostgresqlRepository(conn *sql.DB) domain.FilmsRepository {
	return &filmsPostgresqlRepository{db: conn}
}

func (r *filmsPostgresqlRepository) GetFilmsByGenre(genre string) ([]domain.Film, error) {
	rows, err := r.db.Query(
		`SELECT film.id, film.name, film.preview_path, film.rating FROM film JOIN "genre_film" gf ON film.id = gf.film_id JOIN genre g ON gf.genre_id = g.id WHERE g.name = $1`, genre)
	if err != nil {
		logs.LogError(logs.Logger, "films_postgresql", "GetFilmsByGenre", err, err.Error())
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			logs.LogError(logs.Logger, "postgresql", "GetFilmsByGenre", err, "Failed to close query")
		}
	}(rows)
	logs.Logger.Debug("GetFilmsByGenre query result:", rows)

	var films []domain.Film
	for rows.Next() {
		var film domain.Film
		err = rows.Scan(
			&film.ID,
			&film.Name,
			&film.PreviewPath,
			&film.Rating,
		)

		if err != nil {
			logs.LogError(logs.Logger, "films_postgresql", "GetFilmsByGenre", err, err.Error())
			return nil, err
		}

		films = append(films, film)
	}

	return films, nil
}

func (r *filmsPostgresqlRepository) GetFilmData(id int) (domain.Film, error) {
	row, err := r.db.Query(
		`SELECT *
		FROM film
		WHERE id = $1`, id)
	if err != nil {
		logs.LogError(logs.Logger, "films_postgresql", "GetFilmData", err, err.Error())
		return domain.Film{}, err
	}
	defer func(row *sql.Rows) {
		err := row.Close()
		if err != nil {
			logs.LogError(logs.Logger, "films_postgresql", "GetFilmData", err, err.Error())
		}
	}(row)

	count := 0
	for row.Next() {
		count++
	}

	var film domain.Film
	if count > 0 {
		err = row.Scan(
			&film.ID,
			&film.Name,
			&film.Description,
			&film.Duration,
			&film.PreviewPath,
			&film.MediaPath,
			&film.Country,
			&film.ReleaseYear,
			&film.Rating,
			&film.RatesCount,
			&film.AgeRestriction,
		)
	}

	if err != nil {
		logs.LogError(logs.Logger, "films_postgresql", "GetFilmData", err, err.Error())
		return domain.Film{}, err
	}

	return film, nil
}

func (r *filmsPostgresqlRepository) GetFilmArtists(FilmId int) ([]domain.Artist, error) {
	rows, err := r.db.Query(
		`SELECT name FROM "artist" JOIN "artist-film" af 
		ON id = artist_id 
		WHERE af.film_id = $1`, FilmId)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			logs.LogError(logs.Logger, "films_postgresql", "GetFilmArtists", err, err.Error())
		}
	}(rows)

	var artists []domain.Artist
	for rows.Next() {
		var artist domain.Artist
		err = rows.Scan(
			&artist.Name,
		)
		if err != nil {
			logs.LogError(logs.Logger, "films_postgresql", "GetFilmArtists", err, err.Error())
			return nil, err
		}
		artists = append(artists, artist)
	}

	return artists, nil
}
