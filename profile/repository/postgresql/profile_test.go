package postgres

import (
	"2023_2_Holi/domain"
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/require"
)

const getUserQueryTest = `
	SELECT id, name, email, password, COALESCE\(image_path, ''\) 
	FROM "user" 
`

const updateUserQueryTest = `
	UPDATE "user"
	SET name = \$1, password = \$2, email = \$3, image_path = \$4 
	WHERE id = \$5
	RETURNING id, name, email, password, COALESCE\(image_path, ''\) 
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
		{
			name: "BadCase/EmptyID",
			err:  pgx.ErrNoRows,
		},
		{
			name: "BadCase/IncorrectParam",
			id:   -1,
			err:  errors.New("some error"),
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

func TestUpdateUser(t *testing.T) {
	tests := []struct {
		name    string
		newUser domain.User
		good    bool
		err     error
	}{
		{
			name: "GoodCase/Common",
			newUser: domain.User{
				ID:        1,
				Name:      "Alex",
				Email:     "uvybini@mail.ru",
				Password:  []byte{123},
				ImagePath: "path",
			},
			good: true,
		},
		{
			name:    "GoodCase/EmptyUser",
			newUser: domain.User{},
			good:    true,
		},
		{
			name: "BadCase/NoUser",
			newUser: domain.User{
				ID:        1,
				Name:      "Alex",
				Email:     "uvybini@mail.ru",
				Password:  []byte{123},
				ImagePath: "path",
			},
			err: pgx.ErrNoRows,
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
				AddRow(test.newUser.ID, test.newUser.Name, test.newUser.Email, test.newUser.Password, test.newUser.ImagePath)

			eq := mockDB.ExpectQuery(updateUserQueryTest).
				WithArgs(test.newUser.Name, test.newUser.Password, test.newUser.Email, test.newUser.ImagePath, test.newUser.ID)
			if test.good {
				eq.WillReturnRows(row)
			} else {
				eq.WillReturnError(test.err)
			}

			user, err := r.UpdateUser(test.newUser)
			if test.good {
				require.Nil(t, err)
				require.Equal(t, user, test.newUser)
			} else {
				require.NotNil(t, err)
			}

			err = mockDB.ExpectationsWereMet()
			require.Nil(t, err)
		})
	}
}
