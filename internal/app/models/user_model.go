package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	CommonModel
	OrganizationID *uuid.UUID `gorm:"type:uuid" json:"organization_id"`
	// Organization   organizations.Organization `gorm:"foreignKey:OrganizationID" json:"organization"`
	Roles         string    `gorm:"not null;default:user" json:"roles"`
	Name          string    `gorm:"not null" json:"name"`
	Email         string    `gorm:"unique;not null" json:"email"`
	EmailVerified bool      `gorm:"default:false" json:"email_verified"`
	Image         string    `json:"image,omitempty"`
	IsActive      bool      `gorm:"default:true" json:"is_active"`
	Accounts      []Account `gorm:"foreignKey:UserID" json:"accounts,omitempty"`
}

func (User) TableName() string {
	return "auth.users"
}

// DTOS
type CreateUserDTO struct {
	Name     string `json:"name" validate:"required,min=3"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Roles    string `json:"roles"`
}

type UserResponseDTO struct {
	ID             uuid.UUID  `json:"id"`
	OrganizationID *uuid.UUID `json:"organization_id"`
	Name           string     `json:"name"`
	Email          string     `json:"email"`
	EmailVerified  bool       `json:"email_verified"`
	Roles          string     `json:"roles"`
	Image          string     `json:"image,omitempty"`
	IsActive       bool       `gorm:"default:true" json:"is_active"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}
