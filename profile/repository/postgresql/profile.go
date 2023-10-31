package profile_postgres

import (
	"context"
	"github.com/jackc/pgx/v5"

	"2023_2_Holi/domain"
	logs "2023_2_Holi/logger"

	"github.com/jackc/pgx/v5/pgxpool"
)

const getUserQuery = `
	SELECT id, name, email, password, COALESCE(image_path, '') 
	FROM "user" 
	WHERE id = $1
`

const updateUserQuery = `
	UPDATE "user"
	SET name = $1, password = $2, email = $3, image_path = $4 
	WHERE id = $5 
	RETURNING id, name, email, password, COALESCE(image_path, '') 
`

type profilePostgresqlRepository struct {
	db  *pgxpool.Pool
	ctx context.Context
}

func NewProfilePostgresqlRepository(pool *pgxpool.Pool, ctx context.Context) domain.ProfileRepository {
	return &profilePostgresqlRepository{
		db:  pool,
		ctx: ctx,
	}
}

func (r *profilePostgresqlRepository) GetUser(userID int) (domain.User, error) {
	row := r.db.QueryRow(r.ctx, getUserQuery, userID)

	logs.Logger.Debug("GetProfile query result:", row)

	var user domain.User
	err := row.Scan(
		&user.ID,
		&user.Password,
		&user.Name,
		&user.Email,
		&user.ImagePath,
	)

	if err == pgx.ErrNoRows {
		logs.LogError(logs.Logger, "profile_postgres", "GetUser", err, err.Error())
		return domain.User{}, domain.ErrNotFound
	}
	if err != nil {
		logs.LogError(logs.Logger, "profile_postgres", "GetUser", err, err.Error())
		return domain.User{}, err
	}
	return user, nil
}

func (r *profilePostgresqlRepository) UpdateUser(newUser domain.User) (domain.User, error) {
	row := r.db.QueryRow(r.ctx, updateUserQuery, newUser.Name, newUser.Password, newUser.Email, newUser.ImagePath, newUser.ID)

	var user domain.User
	err := row.Scan(
		&user.ID,
		&user.Password,
		&user.Name,
		&user.Email,
		&user.ImagePath,
	)

	if err == pgx.ErrNoRows {
		logs.LogError(logs.Logger, "profile_postgres", "UpdateProfile", err, err.Error())
		return domain.User{}, domain.ErrNotFound
	}
	if err != nil {
		logs.LogError(logs.Logger, "profile_postgres", "UpdateProfile", err, err.Error())
		return domain.User{}, err
	}
	return user, nil
}
