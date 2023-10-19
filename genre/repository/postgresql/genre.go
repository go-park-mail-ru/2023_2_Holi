package genre_postgres

import (
	"2023_2_Holi/domain"
	logs "2023_2_Holi/logger"
	"database/sql"
)

var logger = logs.LoggerInit()

type genrePostgresRepo struct {
	db *sql.DB
}

func GenrePostgresqlRepository(conn *sql.DB) domain.GenreRepository {
	return &genrePostgresRepo{db: conn}
}

func (r *genrePostgresRepo) GetGenres() ([]domain.Genre, error) {
	rows, err := r.db.Query(
		`SELECT name
		FROM genre`)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			logs.LogError(logger, "postgresql", "GetGenres", err, "Failed to close query")
		}
	}(rows)
	logger.Debug("GetFilmsByGenre query result:", rows)

	var genres []domain.Genre
	for rows.Next() {
		var genre domain.Genre
		err = rows.Scan(
			&genre.Name,
		)

		if err != nil {
			return nil, err
		}

		genres = append(genres, genre)
	}

	return genres, nil
}
