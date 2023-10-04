//package postgresql_test
//
//import (
//	"testing"
//
//	"2023_2_Holi/auth/repository/postgresql"
//
//	_ "github.com/lib/pq"
//	"github.com/stretchr/testify/assert"
//	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
//)
//
//func TestGetByEmail(t *testing.T) {
//	tests := []struct {
//		name     string
//		id       int
//		email    string
//		password string
//		good     bool
//	}{
//		{
//			name:     "GoodCase/Common",
//			id:       1,
//			email:    "fnreo@yandex.ru",
//			password: "txcfygvuhbijn",
//			good:     true,
//		},
//		{
//			name:     "BadCase/EmptyEmail",
//			id:       1,
//			email:    "",
//			password: "txcfygvuhbijn",
//		},
//	}
//
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//			//t.Parallel()
//
//			db, mock, err := sqlmock.New()
//
//			if err != nil {
//				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
//			}
//
//			rows := sqlmock.NewRows([]string{"id", "email", "password"}).
//				AddRow(test.id, test.email, test.password)
//
//			query := `SELECT id, email, password FROM "user" WHERE email = $1`
//
//			mock.ExpectQuery(query).WithArgs(test.email).WillReturnRows(rows)
//			r := postgresql.NewAuthPostgresqlRepository(db)
//
//			user, err := r.GetByEmail(test.email)
//
//			if test.good {
//				assert.NoError(t, err)
//				assert.NotEmpty(t, user)
//			} else {
//				assert.Error(t, err)
//				assert.Empty(t, user)
//			}
//		})
//	}
//}
