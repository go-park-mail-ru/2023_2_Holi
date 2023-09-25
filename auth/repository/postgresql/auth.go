package postgresql

import (
	"database/sql"
	"fmt"

	"2023_2_Holi/domain"
)

type authPostgresqlRepository struct {
	db *sql.DB
}

func NewAuthPostgresqlRepository(conn *sql.DB) domain.AuthRepository {
	return &authPostgresqlRepository{conn}
}

func (r *authPostgresqlRepository) GetByName(name string) (domain.User, error) {
	result, err := r.db.Query(`SELECT name, password FROM "user" WHERE name = $1`, name)
	if err != nil {
		return domain.User{}, err
	}
	defer func(result *sql.Rows) {
		err := result.Close()
		if err != nil {
			fmt.Println("Error: ", err)
		}
	}(result)

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

func (r *authPostgresqlRepository) AddUser(user domain.User) (int, error) {
	if user.Name == "" || user.Password == "" {
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

	var id int
	if err := result.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

type sessionPostgresqlRepository struct {
	db *sql.DB
}

func NewSessionPostgresqlRepository(conn *sql.DB) domain.SessionRepository {
	return &sessionPostgresqlRepository{conn}
}

func (s *sessionPostgresqlRepository) Add(session domain.Session) error {
	_, err := s.db.Exec(`INSERT INTO "session" VALUES ($1, $2, $3)`, session.Token, session.UserID, session.ExpiresAt)
	if err != nil {
		return err
	}
	return nil
}

func (s *sessionPostgresqlRepository) DeleteByToken(token string) error {
	_, err := s.db.Exec(`DELETE FROM "session" WHERE token = $1`, token)
	if err != nil {
		return err
	}
	return nil
}
