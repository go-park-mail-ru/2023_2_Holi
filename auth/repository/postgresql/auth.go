package postgresql

import (
	"2023_2_Holi/domain"
	"database/sql"
)

type userPostgresqlRepository struct {
	db *sql.DB
}

func NewUserPostgresqlRepository(conn *sql.DB) domain.UserRepository {
	return &userPostgresqlRepository{conn}
}

func (r *userPostgresqlRepository) GetByName(name string) (domain.User, error) {
	result, err := r.db.Query(`SELECT name, password FROM "user" WHERE name = $1`, name)
	if err != nil {
		return domain.User{}, err
	}
	defer result.Close()

	var user domain.User
	for result.Next() {
		err = result.Scan(
			&user.Name,
			&user.Password,
		)

		if err != nil {
			return domain.User{}, err
		}
	}

	return user, nil
}
