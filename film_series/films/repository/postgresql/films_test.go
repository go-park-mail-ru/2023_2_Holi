//package postgres
//
//import (
//	"2023_2_Holi/domain"
//	"context"
//	"testing"
//
//	"github.com/jackc/pgx/v5/pgtype"
//	"github.com/pashagolub/pgxmock/v3"
//	"github.com/stretchr/testify/require"
//)
//
//const getByIDFilmData = `
//	SELECT e.name, e.description, e.duration,
//		e.preview_path, e.media_path, preview_video_path, release_year, rating, age_restriction
//	FROM video
//		JOIN episode AS e ON video.id = video_id
//`
//
//const getFilmsByGenreQueryTest = `SELECT DISTINCT v.id, e.name, e.preview_path, v.rating, v.preview_video_path
//FROM video AS v
//JOIN video_cast AS vc ON v.id = vc.video_id
//JOIN "cast" AS c ON vc.cast_id = c.id
//JOIN episode AS e ON e.video_id = v.id
//JOIN video_genre AS vg ON v.id = vg.video_id
//JOIN genre AS g ON vg.genre_id = g.id
//WHERE g.name = \$1\ AND video.seasons_count = 0;
//`
//
//func TestGetFilmsByGenre(t *testing.T) {
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
//					Name:             "Film1",
//					PreviewPath:      "/path/to/preview1",
//					Rating:           8.0,
//					PreviewVideoPath: "/path/to/preview/video1",
//				},
//				{
//					ID:               2,
//					Name:             "Film2",
//					PreviewPath:      "/path/to/preview2",
//					Rating:           7.5,
//					PreviewVideoPath: "/path/to/preview/video2",
//				},
//			},
//			good: true,
//		},
//		// {
//		// 	name:  "GoodCase/EmptyResult",
//		// 	genre: "Comedy",
//		// 	films: []domain.Video{},
//		// 	good:  true,
//		// },
//	}
//
//	mockDB, err := pgxmock.NewPool()
//	if err != nil {
//		t.Fatal(err)
//	}
//	defer mockDB.Close()
//	r := NewFilmsPostgresqlRepository(mockDB, context.Background())
//
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//
//			rows := mockDB.NewRows([]string{"id", "name", "preview_path", "rating", "preview_video_path"})
//
//			for _, film := range test.films {
//				rows.AddRow(film.ID, film.Name, film.PreviewPath, film.Rating, film.PreviewVideoPath)
//			}
//
//			eq := mockDB.ExpectQuery(getFilmsByGenreQueryTest).WithArgs(test.genre)
//
//			if test.good {
//				eq.WillReturnRows(rows)
//			} else {
//				eq.WillReturnError(test.err)
//			}
//
//			films, err := r.GetFilmsByGenre(test.genre)
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
//
//func TestGetFilmData(t *testing.T) {
//	tests := []struct {
//		name string
//		id   int
//		film domain.Video
//		good bool
//		err  error
//	}{
//		{
//			name: "GoodCase/Common",
//			id:   1,
//			film: domain.Video{
//				Name:             "Video Name",
//				Description:      "Video Description",
//				Duration:         pgtype.Interval{},
//				PreviewPath:      "/path/to/preview",
//				MediaPath:        "/path/to/media",
//				PreviewVideoPath: "/path/to/preview/video",
//				ReleaseYear:      2021,
//				Rating:           8.5,
//				AgeRestriction:   16,
//			},
//			good: true,
//		},
//		// {
//		// 	name: "BadCase/NegativeID",
//		// 	id:   -1,
//		// 	err:  pgx.ErrNoRows,
//		// },
//		// {
//		// 	name: "BadCase/NoFilm",
//		// 	id:   1,
//		// 	film: domain.Video{
//		// 		ID:   10,
//		// 		Name: "Avatar",
//		// 	},
//		// 	err: pgx.ErrNoRows,
//		// },
//	}
//
//	mockDB, err := pgxmock.NewPool()
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	if err != nil {
//		t.Fatal(err)
//	}
//	defer mockDB.Close()
//	r := NewFilmsPostgresqlRepository(mockDB, context.Background())
//
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//
//			row := mockDB.NewRows([]string{"id", "name", "duration", "preview_path", "media_path",
//				"preview_video_path", "release_year", "rating", "age_restriction"}).
//				AddRow(test.film.Name, test.film.Description, test.film.Duration, test.film.PreviewPath,
//					test.film.MediaPath, test.film.PreviewVideoPath, test.film.ReleaseYear, test.film.Rating,
//					test.film.AgeRestriction)
//
//			eq := mockDB.ExpectQuery(getFilmDataQuery).WithArgs(test.id)
//
//			if test.good {
//				eq.WillReturnRows(row)
//			} else {
//				eq.WillReturnError(test.err)
//			}
//
//			film, err := r.GetFilmData(test.id)
//			if test.good {
//				require.Nil(t, err)
//				require.Equal(t, film, test.film)
//			} else {
//				require.NotNil(t, err)
//			}
//
//			err = mockDB.ExpectationsWereMet()
//			require.Nil(t, err)
//		})
//
//	}
//}
