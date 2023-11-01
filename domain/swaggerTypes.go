package domain

// for swagger auto doc

type UserRequest struct {
	Name      string `json:"name" binding:"required"`
	Password  string `json:"password" binding:"required"`
	Email     string `json:"email" binding:"required"`
	ImagePath string `json:"imagePath"`
}
