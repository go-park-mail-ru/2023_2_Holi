package artist_postgres

import (
	"context"
	"github.com/jackc/pgx/v5"

	"2023_2_Holi/domain"
	logs "2023_2_Holi/logger"

	"github.com/jackc/pgx/v5/pgxpool"
)

var logger = logs.LoggerInit()

const artistPageQuery = `
	SELECT video.id, e.name, e.preview_path, video.rating
	FROM video
			 JOIN video_cast AS vc ON video_id = vc.video_id
			 JOIN cast AS c ON vc.cast_id = c.id
			 JOIN episode AS e ON e.video_id = video.id
	WHERE c.name = $1
`

type ArtistPostgresqlRepository struct {
	db  *pgxpool.Pool
	ctx context.Context
}

func NewArtistPostgresqlRepository(pool *pgxpool.Pool, ctx context.Context) domain.ArtistRepository {
	return &ArtistPostgresqlRepository{
		db:  pool,
		ctx: ctx,
	}
}

func (r *ArtistPostgresqlRepository) GetArtistPage(name string) ([]domain.Film, error) {
	rows, err := r.db.Query(r.ctx, artistPageQuery, name)
	if err == pgx.ErrNoRows {
		logs.LogError(logs.Logger, "artist_postgres", "GetArtistPage", err, err.Error())
		return nil, domain.ErrNotFound
	}
	if err != nil {
		logs.LogError(logs.Logger, "artist_postgres", "GetArtistPage", err, err.Error())
		return nil, err
	}
	defer rows.Close()
	logs.Logger.Debug("GetArtistPage query result:", rows)

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
