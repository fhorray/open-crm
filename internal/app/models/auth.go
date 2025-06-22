package models

import (
	"time"

	"github.com/google/uuid"
)

// Permissions
// type Role string
// type Permission string

// const (
// 	RoleSuperadmin Role = "superadmin"
// 	RoleAdmin      Role = "admin"
// 	RoleUser       Role = "user"
// )

// const (
// 	PermCreateUser Permission = "create_user"
// 	PermDeleteUser Permission = "delete_user"
// 	PermUpdateUser Permission = "update_user"
// 	PermViewUser   Permission = "view_user"
// )

// Register Request DTO
type RegisterRequestDTO struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// Login Request DTO
type LoginRequestDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// Session Model
type Session struct {
	CommonModel
	UserID                uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	RefreshToken          string    `gorm:"type:text;not null;unique" json:"refresh_token"`
	RefreshTokenExpiresAt time.Time `gorm:"not null" json:"refresh_token_expires_at"`
	UserAgent             string    `gorm:"type:text" json:"user_agent,omitempty"`
	IPAddress             string    `gorm:"type:text" json:"ip_address,omitempty"`
}

func (Session) TableName() string {
	return "auth.sessions"
}

// DTOs

// Used to all responses that contains the sessions data
type SessionDTO struct {
	ID                    uuid.UUID `json:"id"`
	UserID                uuid.UUID `json:"user_id"`
	AccessToken           *string   `json:"access_token"`
	RefreshToken          string    `json:"refresh_token"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
	UserAgent             string    `json:"user_agent,omitempty"`
	IPAddress             string    `json:"ip_address,omitempty"`
}

// Default response when user registers, login or get the current session
type GetSessionResponseDTO struct {
	User    *UserResponseDTO `json:"user"`
	Session *SessionDTO      `json:"session"`
}
