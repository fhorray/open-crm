package models

import (
	"time"

	"github.com/google/uuid"
)

// Permissions
type Role string
type Permission string

const (
	RoleSuperadmin Role = "superadmin"
	RoleAdmin      Role = "admin"
	RoleUser       Role = "user"
)

const (
	PermCreateUser Permission = "create_user"
	PermDeleteUser Permission = "delete_user"
	PermUpdateUser Permission = "update_user"
	PermViewUser   Permission = "view_user"
)

// Register request
type RegisterRequestDTO struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginRequestDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// Session
type Session struct {
	CommonModel
	UserID                uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	AccessToken           string    `gorm:"type:text;not null;unique" json:"access_token"`
	RefreshToken          string    `gorm:"type:text;not null;unique" json:"refresh_token"`
	AccessTokenExpiresAt  time.Time `gorm:"not null" json:"access_token_expires_at"`
	RefreshTokenExpiresAt time.Time `gorm:"not null" json:"refresh_token_expires_at"`
}

func (Session) TableName() string {
	return "core.sessions"
}

type GetSessionResponseDTO struct {
	User    *UserResponseDTO `json:"user"`
	Session *Session         `json:"session"`
}
