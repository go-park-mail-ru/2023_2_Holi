package domain

import (
	"time"
)

type User struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Password   string    `json:"password"`
	Email      string    `json:"email"`
	DateJoined time.Time `json:"dateJoined"`
	ImagePath  string    `json:"imagePath"`
}

type UserUsecase interface {
	Login(user User) (Session, error)
	//Register(w http.ResponseWriter, r *http.Request)
	//Logout(w http.ResponseWriter, r *http.Request)

	//Fetch(ctx context.Context, cursor string, num int64) ([]User, string, error)
	//GetByID(ctx context.Context, id int64) (User, error)
	//Update(ctx context.Context, ar *User) error
	//GetByTitle(ctx context.Context, title string) (User, error)
	//Store(context.Context, *User) error
	//Delete(ctx context.Context, id int64) error
}

type UserRepository interface {
	GetByName(name string) (User, error)
	//GetByID(ctx context.Context, id int64) (User, error)

	//Fetch(ctx context.Context, cursor string, num int64) (res []User, nextCursor string, err error)
	//GetByTitle(ctx context.Context, title string) (User, error)
	//Update(ctx context.Context, ar *User) error
	//Store(ctx context.Context, a *User) error
	//Delete(ctx context.Context, id int64) error
}

type Session struct {
	Token       string
	ExpiresAt   time.Time `json:"expiresAt"`
	SessionData string    `json:"sessionData"`
}
