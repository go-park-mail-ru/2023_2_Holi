package postgresql

import (
	"2023_2_Holi/domain"
	logs "2023_2_Holi/logs"
	"database/sql"
)

var logger = logs.LoggerInit()

type ArtistPostgresqlRepository struct {
	db *sql.DB
}

func NewArtistPostgresqlRepository(conn *sql.DB) domain.ArtistRepository {
	return &ArtistPostgresqlRepository{db: conn}
}

func (r *ArtistPostgresqlRepository) GetArtistPage(name, surname string) ([]domain.Film, error) {
	query := `SELECT film.id, film.name, film.preview_path, film.rating
              FROM film
              JOIN "artist-film" af ON film.id = af.film_id
              JOIN artist a ON af.artist_id = a.id
              WHERE a.name = $1 OR a.surname = $2`

	rows, err := r.db.Query(query, name, surname)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			logs.LogError(logger, "postgresql", "GetArtistPage", err, "Failed to close query")
		}
	}(rows)
	logger.Debug("GetArtistPage query result:", rows)

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
