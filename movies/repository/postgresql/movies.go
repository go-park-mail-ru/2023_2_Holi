package movies_postgresql

import (
	"database/sql"

	"2023_2_Holi/domain"
	logs "2023_2_Holi/logger"
)

type moviesPostgresqlRepository struct {
	db *sql.DB
}

func NewFilmPostgresqlRepository(conn *sql.DB) domain.MoviesRepository {
	return &moviesPostgresqlRepository{db: conn}
}

func (r *moviesPostgresqlRepository) GetMoviesByGenre(genre string) ([]domain.Movie, error) {
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
			logs.LogError(logs.Logger, "movies_postgresql", "GetMoviesByGenre", err, err.Error())
		}
	}(rows)

	var movies []domain.Movie
	for rows.Next() {
		var movie domain.Movie
		err = rows.Scan(
			&movie.ID,
			&movie.Name,
			&movie.PreviewPath,
			&movie.Rating,
		)

		if err != nil {
			return nil, err
		}

		movies = append(movies, movie)
	}

	return movies, nil
}
