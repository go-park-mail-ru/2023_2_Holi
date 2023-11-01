package domain

// for swagger auto doc

type UserRequest struct {
	ID        int    `json:"id" binding:"required"`
	Name      string `json:"name" binding:"required"`
	Password  string `json:"password" binding:"required"`
	Email     string `json:"email" binding:"required"`
	ImagePath string `json:"imagePath"`
}
