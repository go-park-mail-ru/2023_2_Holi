package postgres

import (
	"2023_2_Holi/domain"
	"context"
	"testing"

	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/require"
)

const getSuitableFilmsQueryTest = `
	SELECT id, name, preview_path
	FROM "video", plainto_tsquery\(\$1\) q
	WHERE tsv @@ q
	ORDER BY ts_rank\(tsv, q\) DESC
	LIMIT 10;
`

const getSuitableCastQueryTest = `
	SELECT id, name
	FROM "cast", plainto_tsquery\(\$1\) q
	WHERE tsv @@ q
	ORDER BY ts_rank\(tsv, q\) DESC
	LIMIT 10;
`

func TestGetSuitableFilms(t *testing.T) {
	tests := []struct {
		name      string
		searchStr string
		films     []domain.Film
		good      bool
		err       error
	}{
		{
			name:      "GoodCase/Common",
			searchStr: "Las",
			films: []domain.Film{
				{
					ID:          1,
					Name:        "Las Vegas",
					PreviewPath: "path",
				},
				{
					ID:          2,
					Name:        "Las Alamas",
					PreviewPath: "path",
				},
			},
			good: true,
		},
		{
			name:      "BadCase/NoFound",
			searchStr: "Las",
			films:     []domain.Film{},
			err:       domain.ErrNotFound,
		},
	}

	mockDB, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}

	if err != nil {
		t.Fatal(err)
	}
	defer mockDB.Close()
	r := NewSearchPostgresqlRepository(mockDB, context.Background())

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			rows := mockDB.NewRows([]string{"id", "name", "preview_path"})

			for _, film := range test.films {
				rows.AddRow(film.ID, film.Name, film.PreviewPath)
			}

			eq := mockDB.ExpectQuery(getSuitableFilmsQueryTest).
				WithArgs(test.searchStr)
			if test.good {
				eq.WillReturnRows(rows)
			} else {
				eq.WillReturnError(test.err)
			}

			films, err := r.GetSuitableFilms(test.searchStr)
			if test.good {
				require.Nil(t, err)
				require.Equal(t, films, test.films)
			} else {
				require.NotNil(t, err)
			}

			err = mockDB.ExpectationsWereMet()
			require.Nil(t, err)
		})
	}
}

func TestGetSuitableCast(t *testing.T) {
	tests := []struct {
		name      string
		searchStr string
		cast      []domain.Cast
		good      bool
		err       error
	}{
		{
			name:      "GoodCase/Common",
			searchStr: "lee",
			cast: []domain.Cast{
				{
					ID:   1,
					Name: "Bruce lee",
				},
				{
					ID:   2,
					Name: "lee onardo",
				},
			},
			good: true,
		},
		{
			name:      "BadCase/NoFound",
			searchStr: "Las",
			cast:      []domain.Cast{},
			err:       domain.ErrNotFound,
		},
	}

	mockDB, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}

	if err != nil {
		t.Fatal(err)
	}
	defer mockDB.Close()
	r := NewSearchPostgresqlRepository(mockDB, context.Background())

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			rows := mockDB.NewRows([]string{"id", "name"})

			for _, person := range test.cast {
				rows.AddRow(person.ID, person.Name)
			}

			eq := mockDB.ExpectQuery(getSuitableCastQueryTest).
				WithArgs(test.searchStr)
			if test.good {
				eq.WillReturnRows(rows)
			} else {
				eq.WillReturnError(test.err)
			}

			cast, err := r.GetSuitableCast(test.searchStr)
			if test.good {
				require.Nil(t, err)
				require.Equal(t, cast, test.cast)
			} else {
				require.NotNil(t, err)
			}

			err = mockDB.ExpectationsWereMet()
			require.Nil(t, err)
		})
	}
}
