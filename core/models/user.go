package models

type User struct {
	ID       string `gorm:"primaryKey" json:"id"`
	Name     string
	Email    string `gorm:"unique"`
	Password string
}

type CreateUserDTO struct {
	Name     string `json:"name" validate:"required,min=3"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}
