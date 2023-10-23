package profile_postgres

import (
	"2023_2_Holi/domain"
	"database/sql"
)

type profilePostgresqlRepository struct {
	db *sql.DB
}

func NewProfilePostgresqlRepository(conn *sql.DB) domain.ProfileRepository {
	return &profilePostgresqlRepository{db: conn}
}

func (r *profilePostgresqlRepository) GetProfile(userID int) (domain.User, error) {
	return domain.User{}, nil
}

func (r *profilePostgresqlRepository) UpdateProfile(userID int, newUser domain.User) (domain.User, error) {
	return domain.User{}, nil
}

func (r *profilePostgresqlRepository) UploadImage(userID int, image []byte) error {
	return nil
}
