package postgresql

import (
	"2023_2_Holi/domain"
	logs "2023_2_Holi/logs"
	"database/sql"
)

var logger = logs.LoggerInit()

type FilmPostgresqlRepository struct {
	db *sql.DB
}

func NewFilmPostgresqlRepository(conn *sql.DB) domain.FilmRepository {
	return &FilmPostgresqlRepository{db: conn}
}

func (r *FilmPostgresqlRepository) GetFilmsByGenre(genre string) ([]domain.Film, error) {
	rows, err := r.db.Query(
		`SELECT film.id, film.name, film.preview_path, film.rating
		FROM film
		JOIN "genre_film" gf ON film.id = gf.film_id
		JOIN genre g ON gf.genre_id = g.id
		WHERE g.name = $1`, genre)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			logs.LogError(logger, "postgresql", "GetFilmsByGenre", err, "Failed to close query")
		}
	}(rows)
	logger.Debug("GetFilmsByGenre query result:", rows)

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
			return nil, err
		}

		films = append(films, film)
	}

	return films, nil
}
