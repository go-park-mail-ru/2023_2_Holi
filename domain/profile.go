package domain

import "time"

type User struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Password   string    `json:"password"`
	Email      string    `json:"email"`
	DateJoined time.Time `json:"-"`
	ImagePath  string    `json:"imagePath"`
	ImageData  []byte    `json:"imageData"`
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
