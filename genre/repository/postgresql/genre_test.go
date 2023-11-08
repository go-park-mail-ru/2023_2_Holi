package postgres

import (
	"2023_2_Holi/domain"
	"context"
	"testing"

	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/require"
)

const getGenresQueryTest = `
    SELECT name 
    FROM genre
`

func TestGetGenres(t *testing.T) {
	tests := []struct {
		name   string
		genres []domain.Genre
		good   bool
		err    error
	}{
		{
			name: "GoodCase/Common",
			genres: []domain.Genre{
				{
					Name: "Action",
				},
				{
					Name: "Drama",
				},
				{
					Name: "Comedy",
				},
			},
			good: true,
		},
	}

	mockDB, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mockDB.Close()
	r := GenrePostgresqlRepository(mockDB, context.Background())

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			rows := mockDB.NewRows([]string{"name"})

			for _, genre := range test.genres {
				rows.AddRow(genre.Name)
			}

			eq := mockDB.ExpectQuery(getGenresQueryTest)

			if test.good {
				eq.WillReturnRows(rows)
			} else {
				eq.WillReturnError(test.err)
			}

			genres, err := r.GetGenres()
			if test.good {
				require.Nil(t, err)
				require.Equal(t, test.genres, genres)
			} else {
				require.NotNil(t, err)
			}

			err = mockDB.ExpectationsWereMet()
			require.Nil(t, err)
		})
	}
}
