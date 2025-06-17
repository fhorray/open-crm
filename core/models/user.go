package models

type User struct {
	ID        string  `json:"id" gorm:"primaryKey"`
	Name      string  `json:"name"`
	Email     string  `json:"email" gorm:"unique"`
	Password  string  `json:"-"`
	Image     string  `json:"image"`
	Roles     string  `json:"roles"`
	CompanyID string  `json:"company_id"`
	Company   Company `json:"company" gorm:"foreignKey:CompanyID"`
}

type CreateUserDTO struct {
	Name     string `json:"name" validate:"required,min=3"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Roles    string `json:"roles"`
}
