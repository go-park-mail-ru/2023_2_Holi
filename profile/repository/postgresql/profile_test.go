package postgres

import (
	"2023_2_Holi/domain"
	"context"
	"testing"

	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/require"
)

const getUserQueryTest = `
	SELECT id, name, email, password, COALESCE(image_path, '') 
	FROM "user" 
`

func TestGetUser(t *testing.T) {
	tests := []struct {
		name string
		id   int
		user domain.User
		good bool
		err  error
	}{
		{
			name: "GoodCase/Common",
			id:   1,
			user: domain.User{
				ID:        1,
				Name:      "Alex",
				Email:     "uvybini@mail.ru",
				Password:  []byte{123},
				ImagePath: "path",
			},
			good: true,
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
	r := NewProfilePostgresqlRepository(mockDB, context.Background())

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			row := mockDB.NewRows([]string{"id", "name", "email", "password", "image_path"}).
				AddRow(test.user.ID, test.user.Name, test.user.Email, test.user.Password, test.user.ImagePath)

			eq := mockDB.ExpectQuery(getUserQueryTest).
				WithArgs(test.user.ID)
			if test.good {
				eq.WillReturnRows(row)
			} else {
				eq.WillReturnError(test.err)
			}

			user, err := r.GetUser(test.user.ID)
			if test.good {
				require.Nil(t, err)
				require.Equal(t, user, test.user)
			} else {
				require.NotNil(t, err)
			}

			err = mockDB.ExpectationsWereMet()
			require.Nil(t, err)
		})
	}
}
