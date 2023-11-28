package postgres

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"

	"2023_2_Holi/domain"

	"github.com/pashagolub/pgxmock/v3"
)

const getByEmailQueryTest = `
	SELECT id, email, password
	FROM "user"
`

const addUserQueryTest = `
	INSERT INTO "user"
`

const userExistsQueryTest = `
	SELECT EXISTS\(SELECT 1
				  FROM "user"
				  WHERE email = \$1\)
`

func TestGetByEmail(t *testing.T) {
	tests := []struct {
		name  string
		email string
		user  domain.User
		good  bool
		err   error
	}{
		{
			name:  "GoodCase/Common",
			email: "uvybini@mail.ru",
			user: domain.User{
				ID:       1,
				Email:    "uvybini@mail.ru",
				Password: []byte{123},
			},
			good: true,
		},
		{
			name:  "BadCase/EmptyEmail",
			email: "",
			err:   pgx.ErrNoRows,
		},
		{
			name:  "BadCase/IncorrectParam",
			email: "SELECT",
			err:   errors.New("some error"),
		},
	}

	mockDB, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mockDB.Close()
	r := NewAuthPostgresqlRepository(mockDB, context.Background())

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			row := mockDB.NewRows([]string{"id", "email", "password"}).
				AddRow(test.user.ID, test.user.Email, test.user.Password)

			eq := mockDB.ExpectQuery(getByEmailQueryTest).
				WithArgs(test.user.Email)
			if test.good {
				eq.WillReturnRows(row)
			} else {
				eq.WillReturnError(test.err)
			}

			user, err := r.GetByEmail(test.user.Email)
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

func TestAddUser(t *testing.T) {
	tests := []struct {
		name string
		user domain.User
		good bool
		err  error
	}{
		{
			name: "GoodCase/Common",
			user: domain.User{
				ID:        1,
				Email:     "uvybini@mail.ru",
				Password:  []byte{123},
				ImagePath: "https://fs.kinomania.ru/image/file/film_poster/c/36/c36d4557f6955b864e70fd55a02505f1.228.313.jpeg",
			},
			good: true,
		},
		{
			name: "BadCase/EmptyUser",
			err:  errors.New("some error"),
		},
		{
			name: "BadCase/IncorrectParam",
			user: domain.User{
				ID:        1,
				Email:     "SELECT",
				Password:  []byte{123},
				ImagePath: "https://fs.kinomania.ru/image/file/film_poster/c/36/c36d4557f6955b864e70fd55a02505f1.228.313.jpeg",
			},
			err: errors.New("some error"),
		},
	}

	mockDB, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mockDB.Close()
	r := NewAuthPostgresqlRepository(mockDB, context.Background())

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			row := mockDB.NewRows([]string{"id"}).
				AddRow(test.user.ID)
			eq := mockDB.ExpectQuery(addUserQueryTest).
				WithArgs(test.user.Password, test.user.Name, test.user.Email, test.user.ImagePath)
			if test.good {
				eq.WillReturnRows(row)
			} else {
				eq.Maybe().WillReturnError(test.err)
			}

			id, err := r.AddUser(test.user)
			if test.good {
				require.Nil(t, err)
				require.Equal(t, id, test.user.ID)
			} else {
				require.NotNil(t, err)
			}

			err = mockDB.ExpectationsWereMet()
			require.Nil(t, err)
		})
	}
}

func TestUserExists(t *testing.T) {
	tests := []struct {
		name   string
		email  string
		good   bool
		exists bool
		err    error
	}{
		{
			name:   "GoodCase/Common",
			email:  "uvybini@mail.ru",
			good:   true,
			exists: true,
		},
		{
			name:  "BadCase/EmptyEmail",
			email: "",
			good:  true,
		},
		{
			name:  "BadCase/IncorrectParam",
			email: "SELECT",
			err:   errors.New("some error"),
		},
	}

	mockDB, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mockDB.Close()
	r := NewAuthPostgresqlRepository(mockDB, context.Background())

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			row := mockDB.NewRows([]string{"exists"}).
				AddRow(test.exists)

			eq := mockDB.ExpectQuery(userExistsQueryTest).
				WithArgs(test.email)
			if test.good {
				eq.WillReturnRows(row)
			} else {
				eq.WillReturnError(test.err)
			}

			exists, err := r.UserExists(test.email)
			if test.good {
				require.Nil(t, err)
				require.Equal(t, exists, test.exists)
			} else {
				require.NotNil(t, err)
			}

			err = mockDB.ExpectationsWereMet()
			require.Nil(t, err)
		})
	}
}
