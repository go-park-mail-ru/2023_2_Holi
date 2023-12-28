package postgres

import (
	"2023_2_Holi/domain"
	"context"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/require"
	"testing"
)

const getSeriesByGenreQueryTest = `
    SELECT DISTINCT v.id, v.name, v.preview_path, v.rating, v.preview_video_path, v.seasons_count
    FROM video AS v
        JOIN video_cast AS vc ON v.id = vc.video_id
        JOIN "cast" AS c ON vc.cast_id = c.id
        JOIN episode AS e ON e.video_id = v.id
        JOIN video_genre AS vg ON v.id = vg.video_id
        JOIN genre AS g ON vg.genre_id = g.id
    WHERE g.name = \$1 AND v.seasons_count > 0;
`

const getSeriesDataQueryTest = `
    SELECT DISTINCT video.name, video.description, video.preview_path, preview_video_path, release_year, rating, age_restriction
    FROM video
        JOIN episode AS e ON video.id = e.video_id
    WHERE video.id = \$1 AND video.seasons_count > 0;
`
const getSeriesCastQueryTest = `
    SELECT c.id, c.name
    FROM "cast" c
        JOIN video_cast vc ON c.id = vc.cast_id
        JOIN video v ON vc.video_id = v.id
    WHERE vc.video_id = \$1 AND v.seasons_count > 0;
`

const getSeriesEpisodesQueryTest = `
    SELECT e.id, e.name, e.description, e.media_path, e.number, e.season_number
    FROM episode AS e
    JOIN video AS v ON e.video_id = v.id
    WHERE v.id = \$1 AND v.seasons_count > 0;
`

const getCastPageQueryTest = `
    SELECT video.id, video.name, video.preview_path, video.rating, video.preview_video_path
    FROM video
    WHERE video.id IN \(
        SELECT vc.video_id
        FROM video_cast AS vc
        WHERE vc.cast_id = $1 AND video.seasons_count > 0
    )\;
`

const getCastNameQueryTest = `
    SELECT "cast".name
    FROM "cast" 
    WHERE "cast".id = \$1;
`

//func TestGetSeriesByGenre(t *testing.T) {
//	tests := []struct {
//		name  string
//		genre string
//		films []domain.Video
//		good  bool
//		err   error
//	}{
//		{
//			name:  "GoodCase/Common",
//			genre: "Action",
//			films: []domain.Video{
//				{
//					ID:               1,
//					Name:             "Series1",
//					PreviewPath:      "/path/to/preview1",
//					Rating:           8.0,
//					PreviewVideoPath: "/path/to/preview/video1",
//					SeasonsCount:     3,
//				},
//				{
//					ID:               2,
//					Name:             "Series2",
//					PreviewPath:      "/path/to/preview2",
//					Rating:           7.5,
//					PreviewVideoPath: "/path/to/preview/video2",
//					SeasonsCount:     2,
//				},
//			},
//			good: true,
//		},
//		{
//			name:  "BadCase/GetSeriesByGenreError",
//			genre: "UnknownGenre",
//			err:   domain.ErrNotFound,
//		},
//	}
//
//	mockDB, err := pgxmock.NewPool()
//	if err != nil {
//		t.Fatal(err)
//	}
//	defer mockDB.Close()
//	r := NewSeriesPostgresqlRepository(mockDB, context.Background())
//
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//
//			rows := mockDB.NewRows([]string{"id", "name", "preview_path", "rating", "preview_video_path", "seasons_count"})
//
//			for _, series := range test.films {
//				rows.AddRow(series.ID, series.Name, series.PreviewPath, series.Rating, series.PreviewVideoPath, series.SeasonsCount)
//			}
//
//			eq := mockDB.ExpectQuery(getSeriesByGenreQueryTest).WithArgs(test.genre)
//
//			if test.good {
//				eq.WillReturnRows(rows)
//			} else {
//				eq.WillReturnError(test.err)
//			}
//
//			films, err := r.GetSeriesByGenre(test.id)
//			if test.good {
//				require.Nil(t, err)
//				require.Len(t, films, len(test.films))
//				require.ElementsMatch(t, films, test.films)
//			} else {
//				require.Equal(t, domain.ErrNotFound, err)
//				require.Empty(t, films)
//			}
//		})
//	}
//}

func TestGetSeriesData(t *testing.T) {
	tests := []struct {
		name   string
		id     int
		series domain.Video
		good   bool
		err    error
	}{
		{
			name: "GoodCase/Common",
			id:   1,
			series: domain.Video{
				Name:             "Series Name",
				Description:      "Series Description",
				PreviewPath:      "/path/to/preview",
				PreviewVideoPath: "/path/to/preview/video",
				ReleaseYear:      2021,
				Rating:           8.5,
				AgeRestriction:   16,
			},
			good: true,
		},
		{
			name: "BadCase/GetSeriesByGenreError",
			id:   2,
			err:  domain.ErrNotFound,
		},
	}

	mockDB, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mockDB.Close()
	r := NewSeriesPostgresqlRepository(mockDB, context.Background())

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			row := mockDB.NewRows([]string{"name", "description", "preview_path", "preview_video_path", "release_year", "rating", "age_restriction"}).
				AddRow(test.series.Name, test.series.Description, test.series.PreviewPath, test.series.PreviewVideoPath, test.series.ReleaseYear, test.series.Rating, test.series.AgeRestriction)

			eq := mockDB.ExpectQuery(getSeriesDataQueryTest).WithArgs(test.id)

			if test.good {
				eq.WillReturnRows(row)
			} else {
				eq.WillReturnError(test.err)
			}

			series, err := r.GetSeriesData(test.id)
			if test.good {
				require.Nil(t, err)
				require.Equal(t, series, test.series)
			} else {
				require.NotNil(t, err)
			}

			err = mockDB.ExpectationsWereMet()
			require.Nil(t, err)
		})

	}
}

func TestGetSeriesCast(t *testing.T) {
	tests := []struct {
		name     string
		SeriesID int
		casts    []domain.Cast
		good     bool
		err      error
	}{
		{
			name:     "GoodCase/Common",
			SeriesID: 1,
			casts: []domain.Cast{
				{
					ID:   1,
					Name: "Actor1",
				},
				{
					ID:   2,
					Name: "Actor2",
				},
			},
			good: true,
		},
		{
			name:     "BadCase/GetSeriesCastError",
			SeriesID: 2,
			err:      domain.ErrNotFound,
		},
	}
	mockDB, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mockDB.Close()
	r := NewSeriesPostgresqlRepository(mockDB, context.Background())

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			rows := mockDB.NewRows([]string{"id", "name"})

			for _, cast := range test.casts {
				rows.AddRow(cast.ID, cast.Name)
			}

			mockDB.ExpectQuery(getSeriesCastQueryTest).WithArgs(test.SeriesID).WillReturnRows(rows)

			casts, err := r.GetSeriesCast(test.SeriesID)

			if test.good {
				require.NoError(t, err)
				require.Len(t, casts, len(test.casts))
				require.ElementsMatch(t, casts, test.casts)
			} else {
				require.Equal(t, domain.ErrNotFound, err)
				require.Empty(t, casts)
			}

			err = mockDB.ExpectationsWereMet()
			require.Nil(t, err)
		})
	}
}

func TestGetSeriesEpisodes(t *testing.T) {
	tests := []struct {
		name     string
		SeriesID int
		episodes []domain.Episode
		good     bool
		err      error
	}{
		{
			name:     "GoodCase/Common",
			SeriesID: 1,
			episodes: []domain.Episode{
				{
					ID:          1,
					Name:        "Episode1",
					Description: "Description1",
					MediaPath:   "/path/to/media1",
					Number:      1,
					Season:      1,
				},
				{
					ID:          2,
					Name:        "Episode2",
					Description: "Description2",
					MediaPath:   "/path/to/media2",
					Number:      2,
					Season:      1,
				},
			},
			good: true,
		},
		{
			name:     "BadCase/GetSeriesEpisodesError",
			SeriesID: 2,
			err:      domain.ErrNotFound,
		},
	}

	mockDB, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mockDB.Close()
	r := NewSeriesPostgresqlRepository(mockDB, context.Background())

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			rows := mockDB.NewRows([]string{"id", "name", "description", "media_path", "number", "season_number"})

			for _, episode := range test.episodes {
				rows.AddRow(episode.ID, episode.Name, episode.Description, episode.MediaPath, episode.Number, episode.Season)
			}

			mockDB.ExpectQuery(getSeriesEpisodesQueryTest).WithArgs(test.SeriesID).WillReturnRows(rows)

			episodes, err := r.GetSeriesEpisodes(test.SeriesID)

			if test.good {
				require.NoError(t, err)
				require.Len(t, episodes, len(test.episodes))
				require.ElementsMatch(t, episodes, test.episodes)
			} else {
				require.Equal(t, domain.ErrNotFound, err)
				require.Empty(t, episodes)
			}

			err = mockDB.ExpectationsWereMet()
			require.Nil(t, err)
		})
	}
}

func TestGetCastPageSeries(t *testing.T) {
	tests := []struct {
		name   string
		CastID int
		series []domain.Video
		good   bool
		err    error
	}{
		//{
		//	name:   "GoodCase/Common",
		//	CastID: 1,
		//	series: []domain.Video{
		//		{
		//			ID:               1,
		//			Name:             "Series1",
		//			PreviewPath:      "/path/to/preview1",
		//			Rating:           8.0,
		//			PreviewVideoPath: "/path/to/preview/video1",
		//		},
		//		{
		//			ID:               2,
		//			Name:             "Series2",
		//			PreviewPath:      "/path/to/preview2",
		//			Rating:           7.5,
		//			PreviewVideoPath: "/path/to/preview/video2",
		//		},
		//	},
		//	good: true,
		//},
		// Add more test cases as needed
	}

	mockDB, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mockDB.Close()
	r := NewSeriesPostgresqlRepository(mockDB, context.Background())

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			rows := mockDB.NewRows([]string{"id", "name", "preview_path", "rating", "preview_video_path"})

			for _, series := range test.series {
				rows.AddRow(series.ID, series.Name, series.PreviewPath, series.Rating, series.PreviewVideoPath)
			}

			mockDB.ExpectQuery(getCastPageQueryTest).WithArgs(test.CastID).WillReturnRows(rows)

			series, err := r.GetCastPageSeries(test.CastID)

			if test.good {
				require.NoError(t, err)
				require.Len(t, series, len(test.series))
				require.ElementsMatch(t, series, test.series)
			} else {
				require.Equal(t, domain.ErrNotFound, err)
				require.Empty(t, series)
			}

			err = mockDB.ExpectationsWereMet()
			require.Nil(t, err)
		})
	}
}

func TestGetCastNameSeries(t *testing.T) {
	tests := []struct {
		name   string
		CastID int
		cast   domain.Cast
		good   bool
		err    error
	}{
		{
			name:   "GoodCase/Common",
			CastID: 1,
			cast: domain.Cast{
				ID:   0,
				Name: "Actor1",
			},
			good: true,
		},
	}

	mockDB, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mockDB.Close()
	r := NewSeriesPostgresqlRepository(mockDB, context.Background())

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			rows := mockDB.NewRows([]string{"name"}).AddRow(test.cast.Name)

			mockDB.ExpectQuery(getCastNameQueryTest).WithArgs(test.CastID).WillReturnRows(rows)

			cast, err := r.GetCastNameSeries(test.CastID)

			if test.good {
				require.NoError(t, err)
				require.Equal(t, test.cast, cast)
			} else {
				require.Equal(t, domain.ErrNotFound, err)
				require.Empty(t, cast)
			}

			err = mockDB.ExpectationsWereMet()
			require.Nil(t, err)
		})
	}
}
