package domain

import (
	"github.com/mailru/easyjson"
	"time"
)

//easyjson:json
type Credentials struct {
	Password easyjson.RawMessage `json:"password"`
	//Password []byte `json:"password"`
	Email string `json:"email"`
}

type AuthUsecase interface {
	Login(credentials Credentials) (Session, int, error)
	Logout(token string) error
	Register(user User) (int, error)
	IsAuth(token string) (string, error)
}

type AuthRepository interface {
	GetByEmail(email string) (User, error)
	AddUser(user User) (int, error)
	UserExists(email string) (bool, error)
}

type Session struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
	UserID    int       `json:"-"`
}

type SessionRepository interface {
	Add(session Session) error
	DeleteByToken(token string) error
	SessionExists(token string) (string, error)
}
