package auth_postgres

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"2023_2_Holi/domain"
	logs "2023_2_Holi/logger"
)

const getByEmailQuery = `
	SELECT id, email, password
	FROM "user"
	WHERE email = $1
`

const addUserQuery = `
	INSERT INTO "user" (password, name, email, date_joined, image_path)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id
`

const userExistsQuery = `
	SELECT EXISTS(SELECT 1
				  FROM "user"
				  WHERE email = $1)
`

type authPostgresqlRepository struct {
	db  *pgxpool.Pool
	ctx context.Context
}

func NewAuthPostgresqlRepository(pool *pgxpool.Pool, ctx context.Context) domain.AuthRepository {
	return &authPostgresqlRepository{
		db:  pool,
		ctx: ctx,
	}
}

func (r *authPostgresqlRepository) GetByEmail(email string) (domain.User, error) {
	result, err := r.db.Query(r.ctx, getByEmailQuery, email)
	if err == pgx.ErrNoRows {
		logs.LogError(logs.Logger, "auth_postgres", "GetByEmail", err, err.Error())
		return domain.User{}, domain.ErrNotFound
	}
	if err != nil {
		logs.LogError(logs.Logger, "auth_postgres", "GetByEmail", err, err.Error())
		return domain.User{}, err
	}
	defer result.Close()
	logs.Logger.Debug("GetByEmail query result:", result)

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

	result := r.db.QueryRow(r.ctx, addUserQuery,
		user.Password,
		user.Name,
		user.Email,
		user.DateJoined,
		user.ImagePath)

	logs.Logger.Debug("AddUser queryRow result:", result)

	var id int
	if err := result.Scan(&id); err != nil {
		logs.LogError(logs.Logger, "auth_postgres", "AddUser", err, err.Error())
		return 0, domain.ErrInternalServerError
	}
	return id, nil
}

func (r *authPostgresqlRepository) UserExists(email string) (bool, error) {
	result := r.db.QueryRow(r.ctx, userExistsQuery, email)

	var exist bool
	if err := result.Scan(&exist); err != nil {
		logs.LogError(logs.Logger, "auth_postgres", "UserExists", err, err.Error())
		return false, err
	}

	return exist, nil
}
