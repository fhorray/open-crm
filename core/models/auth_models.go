package models

import "github.com/golang-jwt/jwt/v5"

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

// Login request
type SignInRequestDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// Login Response
type SignInResponseDTO struct {
	AccessToken  string
	RefreshToken string
}

// Get Session Response
type Session struct {
	Aud   string `json:"aud"` // audience
	Email string `json:"email"`
	Exp   int64  `json:"exp"` // timestamp Unix (epoch)
	ID    string `json:"id"`
	Iss   string `json:"iss"` // issuer
	Name  string `json:"name"`
}

type GetSessionResponseDTO struct {
	Session jwt.MapClaims    `json:"session"`
	User    *UserResponseDTO `json:"user"`
}
