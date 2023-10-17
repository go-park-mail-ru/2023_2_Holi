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
		logs.LogError(logs.Logger, "movies_postgresql", "GetMoviesByGenre", err, err.Error())
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
			logs.LogError(logs.Logger, "movies_postgresql", "GetMoviesByGenre", err, err.Error())
			return nil, err
		}

		movies = append(movies, movie)
	}

	return movies, nil
}

func (r *moviesPostgresqlRepository) GetMovieData(id int) (*domain.Movie, error) {
	row, err := r.db.Query(
		`SELECT *
		FROM film
		WHERE id = $1`, id)
	if err != nil {
		logs.LogError(logs.Logger, "movies_postgresql", "GetMovieData", err, err.Error())
		return nil, err
	}
	defer func(row *sql.Rows) {
		err := row.Close()
		if err != nil {
			logs.LogError(logs.Logger, "movies_postgresql", "GetMovieData", err, err.Error())
		}
	}(row)

	movie := new(domain.Movie)
	err = row.Scan(
		&movie.ID,
		&movie.Name,
		&movie.Description,
		&movie.Duration,
		&movie.PreviewPath,
		&movie.MediaPath,
		&movie.Country,
		&movie.ReleaseYear,
		&movie.Rating,
		&movie.RatesCount,
		&movie.AgeRestriction,
	)

	if err != nil {
		logs.LogError(logs.Logger, "movies_postgresql", "GetMovieData", err, err.Error())
		return nil, err
	}

	return movie, nil
}

func (r *moviesPostgresqlRepository) GetMovieArtists(movieId int) ([]domain.Artist, error) {
	rows, err := r.db.Query(
		`SELECT artist_id
		FROM artist-film
		WHERE film_id = $1`, movieId)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			logs.LogError(logs.Logger, "movies_postgresql", "GetMovieArtists", err, err.Error())
		}
	}(rows)

	var artists []domain.Artist
	for rows.Next() {
		var artistID int
		err = rows.Scan(
			&artistID,
		)
		if err != nil {
			logs.LogError(logs.Logger, "movies_postgresql", "GetMovieArtists", err, err.Error())
			return nil, err
		}
		artist, err := r.getArtistDataByID(artistID)
		if err != nil {
			logs.LogError(logs.Logger, "movies_postgresql", "GetMovieArtists", err, err.Error())
			return nil, err
		}
		artists = append(artists, *artist)
	}

	return artists, nil
}

func (r *moviesPostgresqlRepository) getArtistDataByID(id int) (*domain.Artist, error) {
	row, err := r.db.Query(
		`SELECT name
		FROM artist
		WHERE id = $1`, id)
	if err != nil {
		logs.LogError(logs.Logger, "movies_postgresql", "getArtistDataByID", err, err.Error())
		return nil, err
	}
	defer func(row *sql.Rows) {
		err := row.Close()
		if err != nil {
			logs.LogError(logs.Logger, "movies_postgresql", "getArtistDataByID", err, err.Error())
		}
	}(row)

	artist := new(domain.Artist)
	err = row.Scan(
		&artist.ID,
		&artist.Name,
	)

	if err != nil {
		logs.LogError(logs.Logger, "movies_postgresql", "getArtistDataByID", err, err.Error())
		return nil, err
	}

	return artist, nil
}
