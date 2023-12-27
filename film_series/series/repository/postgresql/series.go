package postgres

import (
	"context"

	"2023_2_Holi/domain"
	logs "2023_2_Holi/logger"

	"github.com/jackc/pgx"
)

const getSeriesByGenreQuery = `
    SELECT DISTINCT v.id, v.name, v.preview_path, v.rating, v.preview_video_path, v.seasons_count
    FROM video AS v
        JOIN video_cast AS vc ON v.id = vc.video_id
        JOIN "cast" AS c ON vc.cast_id = c.id
        JOIN episode AS e ON e.video_id = v.id
        JOIN video_genre AS vg ON v.id = vg.video_id
        JOIN genre AS g ON vg.genre_id = g.id
    WHERE g.id = $1 AND v.seasons_count > 0;
`

const getSeriesDataQuery = `
    SELECT DISTINCT video.name, video.description, video.preview_path, preview_video_path, release_year, rating, age_restriction
    FROM video
        JOIN episode AS e ON video.id = e.video_id
    WHERE video.id = $1 AND video.seasons_count > 0;
`

const getSeriesCastQuery = `
    SELECT c.id, c.name
    FROM "cast" c
        JOIN video_cast vc ON c.id = vc.cast_id
        JOIN video v ON vc.video_id = v.id
    WHERE vc.video_id = $1 AND v.seasons_count > 0;
`

const getSeriesEpisodesQuery = `
    SELECT e.id, e.name, e.description, e.media_path, e.number, e.season_number
    FROM episode AS e
    JOIN video AS v ON e.video_id = v.id
    WHERE v.id = $1 AND v.seasons_count > 0;
`

const getCastPageQuery = `
    SELECT video.id, video.name, video.preview_path, video.rating, video.preview_video_path
    FROM video
    WHERE video.id IN (
        SELECT vc.video_id
        FROM video_cast AS vc
        WHERE vc.cast_id = $1 AND video.seasons_count > 0
    );
`

const getCastNameQuery = `
    SELECT "cast".name, "cast".brithday, "cast".place, "cast".carier, "cast".imgpath
    FROM "cast" 
    WHERE "cast".id = $1;
`

const getTopRateSeriesQuery = `	
    SELECT video.id, video.name, video.description, video.preview_video_path
    FROM video
    JOIN episode AS e ON video.id = e.video_id
    WHERE video.seasons_count > 0
    ORDER BY rating DESC
    LIMIT 1;
`

type seriesPostgresqlRepository struct {
	db  domain.PgxPoolIface
	ctx context.Context
}

func NewSeriesPostgresqlRepository(pool domain.PgxPoolIface, ctx context.Context) domain.SeriesRepository {
	return &seriesPostgresqlRepository{
		db:  pool,
		ctx: ctx,
	}
}

func (r *seriesPostgresqlRepository) GetSeriesByGenre(genre int) ([]domain.Video, error) {
	rows, err := r.db.Query(r.ctx, getSeriesByGenreQuery, genre)
	if err != nil {
		logs.LogError(logs.Logger, "series_postgresql", "GetSeriesByGenre", err, err.Error())
		return nil, err
	}
	defer rows.Close()
	logs.Logger.Debug("GetSeriesByGenre query result:", rows)

	var films []domain.Video

	for rows.Next() {
		var film domain.Video
		err = rows.Scan(
			&film.ID,
			&film.Name,
			&film.PreviewPath,
			&film.Rating,
			&film.PreviewVideoPath,
			&film.SeasonsCount,
		)

		if err != nil {
			return nil, err
		}

		films = append(films, film)
	}

	logs.Logger.Info("lenth", len(films))
	if len(films) == 0 {
		return nil, domain.ErrNotFound
	}

	return films, nil
}

func (r *seriesPostgresqlRepository) GetSeriesData(id int) (domain.Video, error) {
	row := r.db.QueryRow(r.ctx, getSeriesDataQuery, id)

	logs.Logger.Debug("GetSeriesData query result:", row)

	var series domain.Video
	err := row.Scan(
		&series.Name,
		&series.Description,
		&series.PreviewPath,
		&series.PreviewVideoPath,
		&series.ReleaseYear,
		&series.Rating,
		&series.AgeRestriction,
	)

	if err == pgx.ErrNoRows {
		logs.LogError(logs.Logger, "series_postgresql", "GetSeriesData", err, "No rows")
		return domain.Video{}, domain.ErrNotFound
	}
	if err != nil {
		logs.LogError(logs.Logger, "series_postgresql", "GetSeriesData", err, err.Error())
		return domain.Video{}, err
	}

	return series, nil
}

func (r *seriesPostgresqlRepository) GetSeriesCast(SeriesId int) ([]domain.Cast, error) {
	rows, err := r.db.Query(r.ctx, getSeriesCastQuery, SeriesId)
	if err != nil {
		logs.LogError(logs.Logger, "series_postgresql", "GetSeriesCast", err, err.Error())
		return nil, err
	}
	defer rows.Close()
	logs.Logger.Debug("GetSeriesCast query result:", rows)

	var artists []domain.Cast
	for rows.Next() {
		var artist domain.Cast
		err = rows.Scan(
			&artist.ID,
			&artist.Name,
		)
		if err != nil {
			logs.LogError(logs.Logger, "series_postgresql", "GetSeriesCast", err, err.Error())
			return nil, err
		}
		artists = append(artists, artist)
	}
	if len(artists) == 0 {
		return nil, domain.ErrNotFound
	}

	return artists, nil
}

func (r *seriesPostgresqlRepository) GetSeriesEpisodes(SeriesId int) ([]domain.Episode, error) {
	rows, err := r.db.Query(r.ctx, getSeriesEpisodesQuery, SeriesId)
	if err != nil {
		logs.LogError(logs.Logger, "series_postgresql", "GetFilmEpisodes", err, err.Error())
		return nil, err
	}
	defer rows.Close()
	logs.Logger.Debug("GetSeriesEpisodes query result:", rows)

	var episodes []domain.Episode
	for rows.Next() {
		var episode domain.Episode
		err = rows.Scan(
			&episode.ID,
			&episode.Name,
			&episode.Description,
			&episode.MediaPath,
			&episode.Number,
			&episode.Season,
		)
		if err != nil {
			logs.LogError(logs.Logger, "series_postgresql", "GetSeriesEpisodes", err, err.Error())
			return nil, err
		}
		episodes = append(episodes, episode)
	}
	if len(episodes) == 0 {
		return nil, domain.ErrNotFound
	}

	return episodes, nil
}

func (r *seriesPostgresqlRepository) GetCastPageSeries(id int) ([]domain.Video, error) {
	rows, err := r.db.Query(r.ctx, getCastPageQuery, id)
	if err != nil {
		logs.LogError(logs.Logger, "series_postgres", "GetCastPageSeries", err, err.Error())
		return nil, err
	}
	defer rows.Close()
	logs.Logger.Debug("GetCastPageSeries query result:", rows)

	var films []domain.Video
	for rows.Next() {
		var film domain.Video
		err = rows.Scan(
			&film.ID,
			&film.Name,
			&film.PreviewPath,
			&film.Rating,
			&film.PreviewVideoPath,
		)

		if err != nil {
			return nil, err
		}

		films = append(films, film)
	}
	if len(films) == 0 {
		return nil, domain.ErrNotFound
	}

	return films, nil
}

func (r *seriesPostgresqlRepository) GetCastNameSeries(id int) (domain.Cast, error) {
	rows := r.db.QueryRow(r.ctx, getCastNameQuery, id)

	logs.Logger.Debug("GetCastNameSeries query result:", rows)

	var cast domain.Cast
	err := rows.Scan(
		&cast.Name,
		&cast.Brithday,
		&cast.Place,
		&cast.Carier,
		&cast.ImgPath,
	)

	if err == pgx.ErrNoRows {
		logs.LogError(logs.Logger, "series_postgresql", "GetCastNameSeries", err, err.Error())
		return domain.Cast{}, domain.ErrNotFound
	}
	if err != nil {
		logs.LogError(logs.Logger, "series_postgresql", "GetCastNameSeries", err, err.Error())
		return domain.Cast{}, err
	}

	return cast, nil
}

func (r *seriesPostgresqlRepository) GetTopRate() (domain.Video, error) {
	rows := r.db.QueryRow(r.ctx, getTopRateSeriesQuery)

	logs.Logger.Debug("GetCastName query result:", rows)

	var video domain.Video
	err := rows.Scan(
		&video.ID,
		&video.Name,
		&video.Description,
		&video.PreviewVideoPath,
	)

	if err == pgx.ErrNoRows {
		logs.LogError(logs.Logger, "series_postgresql", "GetTopRate", err, err.Error())
		return domain.Video{}, domain.ErrNotFound
	}
	if err != nil {
		logs.LogError(logs.Logger, "series_postgresql", "GetTopRate", err, err.Error())
		return domain.Video{}, err
	}

	return video, nil
}
