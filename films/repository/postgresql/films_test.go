package films_postgres

import (
	"2023_2_Holi/domain"
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/require"
)

const getByIDFilmData = `
	SELECT e.name, e.description, e.duration,
		e.preview_path, e.media_path, preview_video_path, release_year, rating, age_restriction
	FROM video
		JOIN episode AS e ON video.id = video_id
`

func TestGetFilmsByGenre(t *testing.T) {
	//TODO Lexa
}

func TestGetFilmData(t *testing.T) {
	tests := []struct {
		name string
		id   int
		film domain.Film
		good bool
		err  error
	}{
		{
			name: "GoodCase/Common",
			id:   1,
			film: domain.Film{
				Name:             "Film Name",
				Description:      "Film Description",
				Duration:         pgtype.Interval{},
				PreviewPath:      "/path/to/preview",
				MediaPath:        "/path/to/media",
				PreviewVideoPath: "/path/to/preview/video",
				ReleaseYear:      2021,
				Rating:           8.5,
				AgeRestriction:   16,
			},
			good: true,
		},
		// {
		// 	name: "BadCase/NegativeID",
		// 	id:   -1,
		// 	err:  pgx.ErrNoRows,
		// },
		// {
		// 	name: "BadCase/NoFilm",
		// 	id:   1,
		// 	film: domain.Film{
		// 		ID:   10,
		// 		Name: "Avatar",
		// 	},
		// 	err: pgx.ErrNoRows,
		// },
	}

	mockDB, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}

	if err != nil {
		t.Fatal(err)
	}
	defer mockDB.Close()
	r := NewFilmsPostgresqlRepository(mockDB, context.Background())

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			row := mockDB.NewRows([]string{"id", "name", "duration", "preview_path", "media_path",
				"preview_video_path", "release_year", "rating", "age_restriction"}).
				AddRow(test.film.Name, test.film.Description, test.film.Duration, test.film.PreviewPath,
					test.film.MediaPath, test.film.PreviewVideoPath, test.film.ReleaseYear, test.film.Rating,
					test.film.AgeRestriction)

			eq := mockDB.ExpectQuery(getFilmDataQuery).WithArgs(test.id)

			if test.good {
				eq.WillReturnRows(row)
			} else {
				eq.WillReturnError(test.err)
			}

			film, err := r.GetFilmData(test.id)
			if test.good {
				require.Nil(t, err)
				require.Equal(t, film, test.film)
			} else {
				require.NotNil(t, err)
			}

			err = mockDB.ExpectationsWereMet()
			require.Nil(t, err)
		})

	}
}

func TestGetFilmCast(t *testing.T) {

}
