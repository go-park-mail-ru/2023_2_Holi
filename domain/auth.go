package domain

import (
	"time"
)

type User struct {
	ID         int       `json:"-"`
	Name       string    `json:"name"`
	Password   string    `json:"password"`
	Email      string    `json:"email"`
	DateJoined time.Time `json:"-"`
	ImagePath  string    `json:"imagePath"`
}

type AuthUsecase interface {
	Login(user User) (Session, error)
	Logout(token string) error
	Register(user User) error
	//Logout(w http.ResponseWriter, r *http.Request)

	//Fetch(ctx context.Context, cursor string, num int64) ([]User, string, error)
	//GetByID(ctx context.Context, id int64) (User, error)
	//Update(ctx context.Context, ar *User) error
	//GetByTitle(ctx context.Context, title string) (User, error)
	//Store(context.Context, *User) error
	//Delete(ctx context.Context, id int64) error
}

type AuthRepository interface {
	GetByName(name string) (User, error)
	AddUser(user User) error
	//GetByID(ctx context.Context, id int64) (User, error)

	//Fetch(ctx context.Context, cursor string, num int64) (res []User, nextCursor string, err error)
	//GetByTitle(ctx context.Context, title string) (User, error)
	//Update(ctx context.Context, ar *User) error
	//Store(ctx context.Context, a *User) error
	//Delete(ctx context.Context, id int64) error
}

type Session struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
	UserID    int       `json:"-"`
}

type SessionRepository interface {
	Add(session Session) error
	DeleteByToken(token string) error
}
