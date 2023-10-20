package auth_postgres

import (
	"database/sql"

	"2023_2_Holi/domain"
	logs "2023_2_Holi/logger"
)

type authPostgresqlRepository struct {
	db *sql.DB
}

func NewAuthPostgresqlRepository(conn *sql.DB) domain.AuthRepository {
	return &authPostgresqlRepository{
		db: conn,
	}
}

func (r *authPostgresqlRepository) GetByEmail(email string) (domain.User, error) {
	result, err := r.db.Query(`SELECT id, email, password FROM "user" WHERE email = $1`, email)
	if err != nil {
		return domain.User{}, err
	}
	defer func(result *sql.Rows) {
		err := result.Close()
		if err != nil {
			logs.LogError(logs.Logger, "postgresql", "GetByName", err, "Failed to close query")
		}
	}(result)
	logs.Logger.Debug("GetByName query result:", result)

	var user domain.User
	for result.Next() {
		err = result.Scan(
			&user.ID,
			&user.Email,
			&user.Password,
		)

		if err != nil {
			return domain.User{}, err
		}
	}

	return user, nil
}

func (r *authPostgresqlRepository) AddUser(user domain.User) (int, error) {
	if user.Email == "" || user.Password == "" {
		return 0, domain.ErrBadRequest
	}

	result := r.db.QueryRow(
		`INSERT INTO "user" (password, name, email, date_joined, image_path) 
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		user.Password,
		user.Name,
		user.Email,
		user.DateJoined,
		user.ImagePath)

	logs.Logger.Debug("AddUser queryRow result:", result)

	var id int
	if err := result.Scan(&id); err != nil {
		return 0, domain.ErrInternalServerError
	}
	return id, nil
}

func (r *authPostgresqlRepository) UserExists(email string) (bool, error) {
	result := r.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM "user" WHERE email = $1)`, email)

	var exist bool
	if err := result.Scan(&exist); err != nil {
		return false, err
	}

	return exist, nil
}
