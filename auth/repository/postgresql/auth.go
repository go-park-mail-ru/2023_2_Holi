package postgres

import (
	"2023_2_Holi/domain"
	logs "2023_2_Holi/logger"
	"context"
	"fmt"

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

const addSubQuery = `
    INSERT INTO "user_subscription" (user_id)
	VALUES ($1)
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

//rows, err := r.db.Query(r.ctx, getFilmsByGenreQuery, genre)
//	if err != nil {
//		logs.LogError(logs.Logger, "films_postgresql", "GetFilmsByGenre", err, err.Error())
//		return nil, err
//	}
//	defer rows.Close()
//	logs.Logger.Debug("GetFilmsByGenre query result:", rows)
//
//	var films []domain.Video
//
//	for rows.Next() {
//		var film domain.Video
//		err = rows.Scan(
//			&film.ID,
//			&film.Name,
//			&film.PreviewPath,
//			&film.Rating,
//			&film.PreviewVideoPath,
//			&film.SeasonsCount,
//		)
//
//		if err != nil {
//			return nil, err
//		}
//
//		films = append(films, film)
//	}
//
//	logs.Logger.Info("lenth", len(films))
//	if len(films) == 0 {
//		return nil, domain.ErrNotFound
//	}
//
//	return films, nil

func (r *authPostgresqlRepository) GetByEmail(email string) (domain.User, error) {
	fmt.Println(21)
	result, err := r.db.Query(r.ctx, getByEmailQuery, email)
	if err != nil {
		logs.LogError(logs.Logger, "postgresql", "GetByEmail", err, err.Error())
		return domain.User{}, err
	}
	fmt.Println(22)
	defer result.Close()
	logs.Logger.Debug("GetByEmail query result:", result)

	var user domain.User
	var c int
	for result.Next() {
		err = result.Scan(
			&user.ID,
			&user.Email,
			&user.Password,
		)
		c++
	}

	if err == pgx.ErrNoRows || c == 0 {
		logs.LogError(logs.Logger, "auth_postgres", "GetByEmail", err, err.Error())
		return domain.User{}, domain.ErrNotFound
	}
	if err != nil {
		logs.LogError(logs.Logger, "auth_postgres", "GetByEmail", err, err.Error())
		return domain.User{}, err
	}

	return user, nil

	//fmt.Println(21)
	//result := r.db.QueryRow(r.ctx, getByEmailQuery, email)
	//fmt.Println(22)
	//
	//logs.Logger.Debug("GetByEmail query result:", result)
	//
	//var user domain.User
	//err := result.Scan(
	//	&user.ID,
	//	&user.Email,
	//	&user.Password,
	//)
	//
	//if err == pgx.ErrNoRows {
	//	logs.LogError(logs.Logger, "auth_postgres", "GetByEmail", err, err.Error())
	//	return domain.User{}, domain.ErrNotFound
	//}
	//if err != nil {
	//	logs.LogError(logs.Logger, "auth_postgres", "GetByEmail", err, err.Error())
	//	return domain.User{}, err
	//}
	//
	//return user, nil
}

func (r *authPostgresqlRepository) AddUser(user domain.User) (int, error) {
	if user.Email == "" || len(user.Password) == 0 {
		return 0, domain.ErrBadRequest
	}

	result, err := r.db.Query(r.ctx, addUserQuery,
		user.Password,
		user.Name,
		user.Email,
		user.ImagePath)

	if err != nil {
		logs.LogError(logs.Logger, "postgresql", "AddUser", err, err.Error())
		return 0, err
	}
	defer result.Close()
	logs.Logger.Debug("AddUser queryRow result:", result)

	var id int
	for result.Next() {
		if err := result.Scan(&id); err != nil {
			logs.LogError(logs.Logger, "auth_postgres", "AddUser", err, err.Error())
			return 0, err
		}
	}

	_, err = r.db.Exec(r.ctx, addSubQuery, id)
	if err != nil {
		logs.LogError(logs.Logger, "postgresql", "AddUser", err, err.Error())
		return 0, err
	}

	//logs.Logger.Debug("AddSub queryRow result:", row)
	return id, nil

	//if user.Email == "" || len(user.Password) == 0 {
	//	return 0, domain.ErrBadRequest
	//}
	//
	//result := r.db.QueryRow(r.ctx, addUserQuery,
	//	user.Password,
	//	user.Name,
	//	user.Email,
	//	user.ImagePath)
	//
	//logs.Logger.Debug("AddUser queryRow result:", result)
	//
	//var id int
	//if err := result.Scan(&id); err != nil {
	//	logs.LogError(logs.Logger, "auth_postgres", "AddUser", err, err.Error())
	//	return 0, err
	//}
	//result = r.db.QueryRow(r.ctx, addSubQuery, id)
	//logs.Logger.Debug("AddSub queryRow result:", result)
	//return id, nil
}

func (r *authPostgresqlRepository) UserExists(email string) (bool, error) {
	result, err := r.db.Query(context.Background(), userExistsQuery, email)
	if err != nil {
		logs.LogError(logs.Logger, "postgresql", "AddUser", err, err.Error())
		return false, err
	}
	defer result.Close()

	var exist bool
	for result.Next() {
		if err := result.Scan(&exist); err != nil {
			logs.LogError(logs.Logger, "auth_postgres", "UserExists", err, err.Error())
			return false, err
		}
	}
	return exist, nil

	//result := r.db.QueryRow(context.Background(), userExistsQuery, email)
	//
	//var exist bool
	//if err := result.Scan(&exist); err != nil {
	//	logs.LogError(logs.Logger, "auth_postgres", "UserExists", err, err.Error())
	//	return false, err
	//}
	//
	//return exist, nil
}
