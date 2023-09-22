package postgresql

import (
	"2023_2_Holi/domain"
)

type userPostgresqlRepository struct {
}

func NewUserPostgresqlRepository() domain.UserRepository {
	return &userPostgresqlRepository{}
}
