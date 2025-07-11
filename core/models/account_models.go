package models

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID                    uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	UserID                uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	User                  User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	ProviderID            string    `gorm:"not null;default:credential" json:"provider_id"`
	AccessToken           string    `json:"access_token,omitempty"`
	RefreshToken          string    `json:"refresh_token,omitempty"`
	AccessTokenExpiresAt  string    `json:"access_token_expires_at,omitempty"`
	RefreshTokenExpiresAt string    `json:"refresh_token_expires_at,omitempty"`
	IDToken               string    `json:"id_token,omitempty"`
	Scope                 string    `json:"scope,omitempty"`
	Password              string    `json:"password,omitempty"`
	CreatedAt             time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt             time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Account) TableName() string {
	return "core.accounts"
}

// DTOS
type CreateAccountDTO struct {
	UserID     string `json:"user_id" validate:"required,uuid4"`
	ProviderID string `json:"provider_id" validate:"required"`
	Password   string `json:"password,omitempty"`
}

type AccountResponseDTO struct {
	ID                    uuid.UUID `json:"id"`
	UserID                uuid.UUID `json:"user_id"`
	User                  User      `json:"user,omitempty"`
	ProviderID            string    `json:"provider_id"`
	AccessToken           string    `json:"access_token"`
	RefreshToken          string    `json:"refresh_token"`
	AccessTokenExpiresAt  string    `json:"access_token_expires_at"`
	RefreshTokenExpiresAt string    `json:"refresh_token_expires_at"`
	IDToken               string    `json:"id_token"`
	Scope                 string    `json:"scope"`
	Password              string    `json:"-"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}
