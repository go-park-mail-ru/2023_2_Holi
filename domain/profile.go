package domain

import (
	"github.com/mailru/easyjson"
	"github.com/microcosm-cc/bluemonday"
)

//easyjson:json
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	//Password  []byte             `json:"-"`
	Password  easyjson.RawMessage `json:"password"`
	Email     string              `json:"email"`
	ImagePath string              `json:"imagePath"`
	//ImageData []byte              `json:"imageData"`
	ImageData easyjson.RawMessage `json:"imageData"`
}

func SanitizeUser(u User, s *bluemonday.Policy) User {
	u.Name = s.Sanitize(u.Name)
	u.Email = s.Sanitize(u.Email)
	u.ImagePath = s.Sanitize(u.ImagePath)
	return u
}

type ProfileUsecase interface {
	GetUserData(userID int) (User, error)
	UpdateUser(newUser User) (User, error)
	UploadImage(userID int, imageData []byte) (string, error)
}

type ProfileRepository interface {
	GetUser(userID int) (User, error)
	UpdateUser(newUser User) (User, error)
}
