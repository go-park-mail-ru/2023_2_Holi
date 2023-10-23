package domain

import "time"

type User struct {
	ID         int       `json:"-"`
	Name       string    `json:"name"`
	Password   string    `json:"password"`
	Email      string    `json:"email"`
	DateJoined time.Time `json:"-"`
	ImagePath  string    `json:"imagePath"`
}

type ProfileUsecase interface {
	GetProfile(userID int) (User, error)
	UpdateProfile(userID int, newUser User) (User, error)
	UploadImage(userID int, image []byte) error
}

type ProfileRepository interface {
	GetProfile(userID int) (User, error)
	UpdateProfile(userID int, newUser User) (User, error)
	UploadImage(userID int, image []byte) error
}
