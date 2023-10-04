package collections_postgresql

import (
	"2023_2_Holi/domain"
	"database/sql"
)

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

		}
	}(rows)

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
