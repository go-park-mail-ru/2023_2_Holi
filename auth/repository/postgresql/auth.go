package postgresql

import (
	"2023_2_Holi/domain"
	"database/sql"
)

type userPostgresqlRepository struct {
	Conn *sql.DB
}

func NewUserPostgresqlRepository(conn *sql.DB) domain.UserRepository {
	return &userPostgresqlRepository{conn}
}
