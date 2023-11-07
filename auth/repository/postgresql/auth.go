package postgres

import (
	"2023_2_Holi/domain"
	logs "2023_2_Holi/logger"
	"context"

	"github.com/jackc/pgx/v5"
)

const getByEmailQuery = `
	SELECT id, email, password
	FROM "user"
	WHERE email = $1
`

const addUserQuery = `
	INSERT INTO "user" (password, name, email, image_path)
	VALUES ($1, $2, $3, $4)
	RETURNING id
`

const userExistsQuery = `
	SELECT EXISTS(SELECT 1
				  FROM "user"
				  WHERE email = $1)
`

type authPostgresqlRepository struct {
	db  domain.PgxPoolIface
	ctx context.Context
}

func NewAuthPostgresqlRepository(pool domain.PgxPoolIface, ctx context.Context) domain.AuthRepository {
	return &authPostgresqlRepository{
		db:  pool,
		ctx: ctx,
	}
}

func (r *authPostgresqlRepository) GetByEmail(email string) (domain.User, error) {
	result := r.db.QueryRow(r.ctx, getByEmailQuery, email)

	logs.Logger.Debug("GetByEmail query result:", result)

	var user domain.User
	err := result.Scan(
		&user.ID,
		&user.Email,
		&user.Password,
	)

	if err == pgx.ErrNoRows {
		logs.LogError(logs.Logger, "auth_postgres", "GetByEmail", err, err.Error())
		return domain.User{}, domain.ErrNotFound
	}
	if err != nil {
		logs.LogError(logs.Logger, "auth_postgres", "GetByEmail", err, err.Error())
		return domain.User{}, err
	}

	return user, nil
}

func (r *authPostgresqlRepository) AddUser(user domain.User) (int, error) {
	if user.Email == "" || len(user.Password) == 0 {
		return 0, domain.ErrBadRequest
	}

	result := r.db.QueryRow(r.ctx, addUserQuery,
		user.Password,
		user.Name,
		user.Email,
		user.ImagePath)

	logs.Logger.Debug("AddUser queryRow result:", result)

	var id int
	if err := result.Scan(&id); err != nil {
		logs.LogError(logs.Logger, "auth_postgres", "AddUser", err, err.Error())
		return 0, err
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
