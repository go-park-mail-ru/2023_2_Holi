package profile_postgres

import (
	"2023_2_Holi/domain"
	logs "2023_2_Holi/logger"
	"database/sql"
)

type profilePostgresqlRepository struct {
	db *sql.DB
}

func NewProfilePostgresqlRepository(conn *sql.DB) domain.ProfileRepository {
	return &profilePostgresqlRepository{db: conn}
}

func (r *profilePostgresqlRepository) GetUser(userID int) (domain.User, error) {
	row, err := r.db.Query(
		`SELECT * FROM user WHERE id = $1`, userID)
	if err != nil {
		logs.LogError(logs.Logger, "profile_postgres", "GetProfile", err, err.Error())
		return domain.User{}, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			logs.LogError(logs.Logger, "profile_postgres", "GetProfile", err, "Failed to close query")
		}
	}(row)
	logs.Logger.Debug("GetProfile query result:", row)

	var user domain.User
	err = row.Scan(
		&user.ID,
		&user.Password,
		&user.Name,
		&user.Email,
		&user.DateJoined,
		&user.ImagePath,
	)

	return user, nil
}

func (r *profilePostgresqlRepository) UpdateUser(newUser domain.User) (domain.User, error) {
	row, err := r.db.Query(
		`UPDATE users SET name = $1, password = $2, email = $3, image_path = $4 WHERE id = $5 RETURNING *`,
		newUser.Name, newUser.Password, newUser.Email, newUser.ImagePath, newUser.ID)
	if err != nil {
		logs.LogError(logs.Logger, "profile_postgres", "UpdateProfile", err, err.Error())
		return domain.User{}, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			logs.LogError(logs.Logger, "profile_postgres", "UpdateProfile", err, "Failed to close query")
		}
	}(row)
	logs.Logger.Debug("UpdateProfile query result:", row)

	var user domain.User
	err = row.Scan(
		&user.ID,
		&user.Password,
		&user.Name,
		&user.Email,
		&user.DateJoined,
		&user.ImagePath,
	)

	return user, nil
}
