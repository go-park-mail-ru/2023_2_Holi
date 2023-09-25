package domain

import (
	"time"
)

type Credentials struct {
	Password string `json:"password"`
	Name     string `json:"name"`
}

type User struct {
	ID         int       `json:"-"`
	Name       string    `json:"name"`
	Password   string    `json:"password"`
	Email      string    `json:"email"`
	DateJoined time.Time `json:"-"`
	ImagePath  string    `json:"imagePath"`
}

type AuthUsecase interface {
	Login(credentials Credentials) (Session, error)
	Logout(token string) error
	Register(user User) (int, error)
}

type AuthRepository interface {
	GetByName(name string) (User, error)
	AddUser(user User) (int, error)
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
